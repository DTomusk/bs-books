package events

import (
	"bs-books-api/internal/testutil"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublishAndDequeue(t *testing.T) {
	testutil.WithTx(t, func(tx *sql.Tx) {
		// Arrange
		txRunner := testutil.NewTestTxRunner(tx)
		repo := NewEventRepo()
		service := NewEventService(txRunner, repo, 5)
		ctx := context.Background()

		// Act 1: publish event to queue
		err := service.PublishEvent(
			ctx,
			"TestEventType",
			"43681e21-08d4-43e1-b0b6-8d6f75a9b8b3",
			map[string]string{"key": "value"},
		)

		// Assert 1
		require.NoError(t, err)

		// Act 2: dequeue an event from queue
		event, err := service.DequeueEvent(ctx)

		// Assert 2
		require.NoError(t, err)
		require.NotNil(t, event)
		require.Equal(t, "TestEventType", event.Type)
		require.Equal(t, "43681e21-08d4-43e1-b0b6-8d6f75a9b8b3", event.AggregateID)
		require.JSONEq(t, `{"key": "value"}`, string(event.Payload))
	})
}
