package workers

import (
	"sync"
)

// LogMessage represents a single log entry
type LogMessage struct {
	DeploymentID string `json:"deploymentId"`
	Message      string `json:"message"`
	Timestamp    int64  `json:"timestamp"`
}

// LogBroadcaster manages WebSocket subscribers for deployment logs
type LogBroadcaster struct {
	mu          sync.RWMutex
	subscribers map[string][]chan LogMessage // deploymentID -> channels
}

var (
	broadcaster     *LogBroadcaster
	broadcasterOnce sync.Once
)

// GetBroadcaster returns the singleton LogBroadcaster instance
func GetBroadcaster() *LogBroadcaster {
	broadcasterOnce.Do(func() {
		broadcaster = &LogBroadcaster{
			subscribers: make(map[string][]chan LogMessage),
		}
	})
	return broadcaster
}

// Subscribe adds a new subscriber for a deployment
func (b *LogBroadcaster) Subscribe(deploymentID string) chan LogMessage {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan LogMessage, 100) // buffered channel
	b.subscribers[deploymentID] = append(b.subscribers[deploymentID], ch)
	return ch
}

// Unsubscribe removes a subscriber
func (b *LogBroadcaster) Unsubscribe(deploymentID string, ch chan LogMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()

	channels := b.subscribers[deploymentID]
	for i, c := range channels {
		if c == ch {
			// Remove channel from slice
			b.subscribers[deploymentID] = append(channels[:i], channels[i+1:]...)
			close(ch)
			break
		}
	}

	// Clean up empty deployment entries
	if len(b.subscribers[deploymentID]) == 0 {
		delete(b.subscribers, deploymentID)
	}
}

// Broadcast sends a log message to all subscribers of a deployment
func (b *LogBroadcaster) Broadcast(msg LogMessage) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	channels := b.subscribers[msg.DeploymentID]
	for _, ch := range channels {
		select {
		case ch <- msg:
			// Message sent successfully
		default:
			// Channel is full, skip this message (non-blocking)
		}
	}
}

// GetSubscriberCount returns the number of subscribers for a deployment
func (b *LogBroadcaster) GetSubscriberCount(deploymentID string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subscribers[deploymentID])
}
