package routes

import (
	"net/http"

	"deployes/internal/interfaces/http/handlers"
)

func RegisterWebSocketRoutes(handler *handlers.WebSocketHandler) {
	// WebSocket route for deployment log streaming
	// Note: WebSocket connections use token in query params instead of headers
	http.HandleFunc("/api/deployments/", func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a WebSocket log streaming request
		// Expected path: /api/deployments/{id}/logs/stream
		if len(r.URL.Path) > len("/api/deployments/") {
			// Let the WebSocket handler validate and process
			handler.StreamDeploymentLogs(w, r)
		}
	})
}
