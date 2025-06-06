package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/carped99/gosdk/events"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type defaultPublisher struct {
	executor Executor
}

func (p *defaultPublisher) Publish(ctx context.Context, messages *events.Message) error {
	payload, err := json.Marshal(messages.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message payload: %w", err)
	}

	metadata, err := json.Marshal(messages.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal message metadata: %w", err)
	}

	query := `
        INSERT INTO outbox_message (
            uuid, event_topic, event_domain, event_type, object_type, producer, correlation_id, payload, metadata, created_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        )
    `

	_, err = p.executor.ExecContext(ctx, query,
		messages.UUID,
		messages.EventTopic,
		messages.EventDomain,
		messages.EventType,
		messages.ObjectType,
		messages.Producer,
		messages.CorrelationID,
		string(payload),
		string(metadata),
		messages.CreatedAt,
	)

	return err
}

func NewPublisher(executor Executor) events.Publisher {
	return &defaultPublisher{
		executor: executor,
	}
}
