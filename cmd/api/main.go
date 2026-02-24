package main

import (
	"log"
	"net/http"

	"deployes/internal/config"
	"deployes/internal/infrastucture/database/postgres"
	"deployes/internal/infrastucture/repositories"

	appDeployment "deployes/internal/application/deployment"
	appProject "deployes/internal/application/project"
	appSecret "deployes/internal/application/secret"
	appServer "deployes/internal/application/server"
	appUser "deployes/internal/application/user"
	appWebhook "deployes/internal/application/webhook"
	"deployes/internal/infrastucture/workers"

	httpHandlers "deployes/internal/interfaces/http/handlers"
	httpRoutes "deployes/internal/interfaces/http/routes"
)

func main() {
	// 1) Config yükle
	cfg := config.Load()

	// 2) DB bağlantısı
	db := postgres.NewConnection(cfg.DBUrl)
	postgres.RunMigrations(cfg.DBUrl) // Initialize schema
	defer db.Close()

	// 3) Dependencies (Repo -> Service -> Handler)
	userRepo := repositories.NewUserRepository(db)
	userService := appUser.NewService(userRepo, cfg.JWTSecret)
	userHandler := httpHandlers.NewUserHandler(userService)

	projectRepo := repositories.NewProjectRepository(db)
	projectService := appProject.NewService(projectRepo, cfg.EncryptionKey)
	projectHandler := httpHandlers.NewProjectHandler(projectService)

	serverRepo := repositories.NewServerRepository(db)
	serverService := appServer.NewService(serverRepo, cfg.EncryptionKey)
	serverHandler := httpHandlers.NewServerHandler(serverService)

	deploymentRepo := repositories.NewDeploymentRepository(db)
	deploymentService := appDeployment.NewService(deploymentRepo, projectRepo, serverRepo)
	deploymentHandler := httpHandlers.NewDeploymentHandler(deploymentService)

	webhookRepo := repositories.NewWebhookRepository(db)
	webhookService := appWebhook.NewService(webhookRepo, cfg.BaseURL)
	webhookHandler := httpHandlers.NewWebhookHandler(webhookService, deploymentService, projectService)

	websocketHandler := httpHandlers.NewWebSocketHandler(cfg.JWTSecret, cfg.AllowedOrigins)

	encryptionHandler := httpHandlers.NewEncryptionHandler(cfg)

	secretRepo := repositories.NewSecretRepository(db)
	secretService := appSecret.NewService(secretRepo, cfg.EncryptionKey)
	secretHandler := httpHandlers.NewSecretHandler(secretService)

	terminalHandler := httpHandlers.NewTerminalHandler(serverRepo, cfg.EncryptionKey, cfg.JWTSecret, cfg.AllowedOrigins)

	// worker
	deployWorker := workers.NewDeployWorker(
		deploymentRepo,
		projectRepo,
		serverRepo,
		secretRepo,
		cfg.EncryptionKey,
	)
	go deployWorker.Start()

	// 4) Routes
	httpRoutes.RegisterAllRoutes(cfg, userHandler, projectHandler, serverHandler, deploymentHandler, webhookHandler, websocketHandler, encryptionHandler, secretHandler)
	httpRoutes.RegisterTerminalRoutes(cfg, terminalHandler)

	// 5) Health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("deployes API is running ✅"))
	})

	log.Printf("🚀 Server started on port %s", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, nil))
}
