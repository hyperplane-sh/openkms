package audit

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	TYPE_FILE = "file"

	LEVEL_INFO  = "INFO"
	LEVEL_WARN  = "WARN"
	LEVEL_ERROR = "ERROR"

	TOPIC_LIFECYCLE = "LIFECYCLE"
)

type Auditor interface {
	// RecordEvent - records an auditing event.
	RecordEvent(event Event) error
	// Persist - persists any buffered events to the storage.
	Persist() error
	// Close - closes the auditor and releases any resources.
	Close() error
}

type Event struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Group     string            `json:"group"`
	Topic     string            `json:"topic"`
	Message   string            `json:"message"`
	Labels    map[string]string `json:"labels"`
}

func NewEvent(level, group, topic, message string, labels map[string]string) Event {
	return Event{
		Timestamp: time.Now(),
		Level:     level,
		Group:     group,
		Topic:     topic,
		Message:   message,
		Labels:    labels,
	}
}

// ToString - converts the event to a string representation.
func (e Event) ToString() string {
	return fmt.Sprintf("[%s] [%s] [%s] [%s] %s Labels: %v", e.Timestamp.Format(time.RFC3339), e.Level, e.Group, e.Topic, e.Message, e.Labels)
}

// ToJSON - converts the event to a JSON representation.
func (e Event) ToJSON() string {
	content, _ := json.Marshal(e)
	return string(content)
}
