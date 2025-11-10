package websocket

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Client -> Server message types
	MessageTypeSubscribe   MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"
	MessageTypePing        MessageType = "ping"

	// Server -> Client message types
	MessageTypeJobUpdate   MessageType = "job_update"
	MessageTypeJobComplete MessageType = "job_complete"
	MessageTypeError       MessageType = "error"
	MessageTypePong        MessageType = "pong"
)

// ClientMessage represents a message from client to server
type ClientMessage struct {
	Type  MessageType `json:"type"`
	JobID string      `json:"jobId,omitempty"`
}

// ServerMessage represents a message from server to client
type ServerMessage struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}
