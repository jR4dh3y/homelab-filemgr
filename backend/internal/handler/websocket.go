package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/homelab/filemanager/internal/service"
	ws "github.com/homelab/filemanager/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any origin in development
	// In production, this should be restricted
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub         *ws.Hub
	authService service.AuthService
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *ws.Hub, authService service.AuthService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
	}
}

// ServeWS handles WebSocket upgrade requests with authentication
func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Extract token from query parameter or Authorization header
	token := h.extractToken(r)
	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	// Validate the token
	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
		return
	}

	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Upgrade already sends an error response
		return
	}

	// Create a new client and register with the hub
	client := ws.NewClient(h.hub, conn, claims.UserID)
	h.hub.Register(client)

	// Start the client's read and write pumps in separate goroutines
	go client.WritePump()
	go client.ReadPump()
}

// extractToken extracts the JWT token from the request
// It checks the query parameter 'token' first, then the Authorization header
func (h *WebSocketHandler) extractToken(r *http.Request) string {
	// Check query parameter first (useful for WebSocket connections)
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Remove "Bearer " prefix if present
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}
