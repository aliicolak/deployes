package routes

import (
	"net/http"

	"deployes/internal/config"
	"deployes/internal/interfaces/http/handlers"
)

func RegisterTerminalRoutes(cfg *config.Config, handler *handlers.TerminalHandler) {
	// No middleware here because Auth is handled inside via Query Param
	http.HandleFunc("/api/ws/terminal", handler.HandleWebSocket)
}
