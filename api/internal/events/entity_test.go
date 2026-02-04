package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEvent(t *testing.T) {
	// Arrange
	payload := struct {
		Message string `json:"message"`
	}{
		Message: "Test event",
	}

	// Act
	event, err := newEvent("TestEvent", "aggregate-id", payload)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, event)
	require.Equal(t, "TestEvent", event.Type)
	require.Equal(t, "aggregate-id", event.AggregateID)
	require.NotEmpty(t, event.ID)
	require.NotEmpty(t, event.OccurredAt)
}
