package audit_test

import (
	"testing"
	"time"

	"github.com/hyperplane-sh/openkms/internal/audit"
)

func TestNewEvent(t *testing.T) {
	scenarios := []struct {
		name       string
		level      string
		group      string
		topic      string
		message    string
		labels     map[string]string
		assertions func(t *testing.T, event audit.Event)
	}{
		{
			name:    "Basic Event Creation",
			level:   audit.LEVEL_INFO,
			group:   "TestGroup",
			topic:   audit.TOPIC_LIFECYCLE,
			message: "This is a test event",
			labels:  map[string]string{"key1": "value1"},
			assertions: func(t *testing.T, event audit.Event) {
				if event.Level != audit.LEVEL_INFO {
					t.Errorf("Expected level %s, got %s", audit.LEVEL_INFO, event.Level)
				}
				if event.Group != "TestGroup" {
					t.Errorf("Expected group %s, got %s", "TestGroup", event.Group)
				}
				if event.Topic != audit.TOPIC_LIFECYCLE {
					t.Errorf("Expected topic %s, got %s", audit.TOPIC_LIFECYCLE, event.Topic)
				}
				if event.Message != "This is a test event" {
					t.Errorf("Expected message %s, got %s", "This is a test event", event.Message)
				}
				if val, ok := event.Labels["key1"]; !ok || val != "value1" {
					t.Errorf("Expected label key1 to be value1, got %v", event.Labels)
				}
			},
		},
		{
			name:    "Event Creation with Multiple Labels",
			level:   audit.LEVEL_WARN,
			group:   "AnotherGroup",
			topic:   "CUSTOM_TOPIC",
			message: "Warning event",
			labels:  map[string]string{"env": "prod", "version": "1.0"},
			assertions: func(t *testing.T, event audit.Event) {
				if event.Level != audit.LEVEL_WARN {
					t.Errorf("Expected level %s, got %s", audit.LEVEL_WARN, event.Level)
				}
				if event.Group != "AnotherGroup" {
					t.Errorf("Expected group %s, got %s", "AnotherGroup", event.Group)
				}
				if event.Topic != "CUSTOM_TOPIC" {
					t.Errorf("Expected topic %s, got %s", "CUSTOM_TOPIC", event.Topic)
				}
				if event.Message != "Warning event" {
					t.Errorf("Expected message %s, got %s", "Warning event", event.Message)
				}
				if val, ok := event.Labels["env"]; !ok || val != "prod" {
					t.Errorf("Expected label env to be prod, got %v", event.Labels)
				}
				if val, ok := event.Labels["version"]; !ok || val != "1.0" {
					t.Errorf("Expected label version to be 1.0, got %v", event.Labels)
				}
			},
		},
		{
			name:    "Event Creation with No Labels",
			level:   audit.LEVEL_ERROR,
			group:   "ErrorGroup",
			topic:   "ERROR_TOPIC",
			message: "Error occurred",
			labels:  nil,
			assertions: func(t *testing.T, event audit.Event) {
				if event.Level != audit.LEVEL_ERROR {
					t.Errorf("Expected level %s, got %s", audit.LEVEL_ERROR, event.Level)
				}
				if event.Group != "ErrorGroup" {
					t.Errorf("Expected group %s, got %s", "ErrorGroup", event.Group)
				}
				if event.Topic != "ERROR_TOPIC" {
					t.Errorf("Expected topic %s, got %s", "ERROR_TOPIC", event.Topic)
				}
				if event.Message != "Error occurred" {
					t.Errorf("Expected message %s, got %s", "Error occurred", event.Message)
				}
				if event.Labels != nil {
					t.Errorf("Expected labels to be nil, got %v", event.Labels)
				}
			},
		},
		{
			name:    "Event Creation with Empty Labels",
			level:   audit.LEVEL_INFO,
			group:   "EmptyLabelGroup",
			topic:   "EMPTY_LABEL_TOPIC",
			message: "Event with empty labels",
			labels:  map[string]string{},
			assertions: func(t *testing.T, event audit.Event) {
				if event.Level != audit.LEVEL_INFO {
					t.Errorf("Expected level %s, got %s", audit.LEVEL_INFO, event.Level)
				}
				if event.Group != "EmptyLabelGroup" {
					t.Errorf("Expected group %s, got %s", "EmptyLabelGroup", event.Group)
				}
				if event.Topic != "EMPTY_LABEL_TOPIC" {
					t.Errorf("Expected topic %s, got %s", "EMPTY_LABEL_TOPIC", event.Topic)
				}
				if event.Message != "Event with empty labels" {
					t.Errorf("Expected message %s, got %s", "Event with empty labels", event.Message)
				}
				if len(event.Labels) != 0 {
					t.Errorf("Expected labels to be empty, got %v", event.Labels)
				}
			},
		},
		{
			name:    "Event Creation with Special Characters in Message",
			level:   audit.LEVEL_INFO,
			group:   "SpecialCharGroup",
			topic:   "SPECIAL_CHAR_TOPIC",
			message: "Event with special characters !@#$%^&*()",
			labels:  map[string]string{"special": "chars"},
			assertions: func(t *testing.T, event audit.Event) {
				if event.Level != audit.LEVEL_INFO {
					t.Errorf("Expected level %s, got %s", audit.LEVEL_INFO, event.Level)
				}
				if event.Group != "SpecialCharGroup" {
					t.Errorf("Expected group %s, got %s", "SpecialCharGroup", event.Group)
				}
				if event.Topic != "SPECIAL_CHAR_TOPIC" {
					t.Errorf("Expected topic %s, got %s", "SPECIAL_CHAR_TOPIC", event.Topic)
				}
				if event.Message != "Event with special characters !@#$%^&*()" {
					t.Errorf("Expected message %s, got %s", "Event with special characters !@#$%^&*()", event.Message)
				}
				if val, ok := event.Labels["special"]; !ok || val != "chars" {
					t.Errorf("Expected label special to be chars, got %v", event.Labels)
				}
			},
		},
		{
			name:    "Event Creation with Trash Labels",
			level:   audit.LEVEL_INFO,
			group:   "TrashLabelGroup",
			topic:   "TRASH_LABEL_TOPIC",
			message: "Event with trash labels",
			labels:  map[string]string{"": "", "   ": "   "},
			assertions: func(t *testing.T, event audit.Event) {
				if event.Level != audit.LEVEL_INFO {
					t.Errorf("Expected level %s, got %s", audit.LEVEL_INFO, event.Level)
				}
				if event.Group != "TrashLabelGroup" {
					t.Errorf("Expected group %s, got %s", "TrashLabelGroup", event.Group)
				}
				if event.Topic != "TRASH_LABEL_TOPIC" {
					t.Errorf("Expected topic %s, got %s", "TRASH_LABEL_TOPIC", event.Topic)
				}
				if event.Message != "Event with trash labels" {
					t.Errorf("Expected message %s, got %s", "Event with trash labels", event.Message)
				}
				if val, ok := event.Labels[""]; !ok || val != "" {
					t.Errorf("Expected label '' to be '', got %v", event.Labels)
				}
				if val, ok := event.Labels["   "]; !ok || val != "   " {
					t.Errorf("Expected label '   ' to be '   ', got %v", event.Labels)
				}
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			event := audit.NewEvent(scenario.level, scenario.group, scenario.topic, scenario.message, scenario.labels)
			scenario.assertions(t, event)
		})
	}
}

func TestEventToString(t *testing.T) {
	// Test with labels
	//
	event := audit.NewEvent(
		audit.LEVEL_INFO,
		"TestGroup",
		audit.TOPIC_LIFECYCLE,
		"This is a test event",
		map[string]string{"key1": "value1"},
	)

	expectedPrefix := "[" + event.Timestamp.Format(time.RFC3339) + "] [INFO] [TestGroup] [LIFECYCLE] This is a test event Labels: map[key1:value1]"
	actual := event.ToString()
	if actual != expectedPrefix {
		t.Errorf("Expected ToString output to be %s, got %s", expectedPrefix, actual)
	}

	// Test without labels
	//
	eventWithoutLabels := audit.NewEvent(
		audit.LEVEL_WARN,
		"AnotherGroup",
		"CUSTOM_TOPIC",
		"Warning event",
		nil,
	)

	expectedPrefix = "[" + eventWithoutLabels.Timestamp.Format(time.RFC3339) + "] [WARN] [AnotherGroup] [CUSTOM_TOPIC] Warning event Labels: map[]"
	actual = eventWithoutLabels.ToString()
	if actual != expectedPrefix {
		t.Errorf("Expected ToString output to be %s, got %s", expectedPrefix, actual)
	}

	// Test with empty labels
	//
	eventWithEmptyLabels := audit.NewEvent(
		audit.LEVEL_ERROR,
		"ErrorGroup",
		"ERROR_TOPIC",
		"Error occurred",
		map[string]string{},
	)

	expectedPrefix = "[" + eventWithEmptyLabels.Timestamp.Format(time.RFC3339) + "] [ERROR] [ErrorGroup] [ERROR_TOPIC] Error occurred Labels: map[]"
	actual = eventWithEmptyLabels.ToString()
	if actual != expectedPrefix {
		t.Errorf("Expected ToString output to be %s, got %s", expectedPrefix, actual)
	}
}
