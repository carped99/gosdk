package outbox

import (
	"context"
	"database/sql"
	"time"
)

type Publisher struct {
	executor Executor
	config   PublisherConfig
}

type PublisherConfig struct {
	BatchSize  int
	MaxRetries int
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func NewPublisher(executor Executor, config PublisherConfig) *Publisher {
	return &Publisher{
		executor: executor,
		config:   config,
	}
}

func (p *Publisher) Publish(ctx context.Context, messages ...Message) error {
	if len(messages) == 0 {
		return nil
	}
	return p.publishBatch(ctx, messages)
}

func (p *Publisher) publishBatch(ctx context.Context, messages []Message) error {
	for _, msg := range messages {
		if err := p.publishWithRetry(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func (p *Publisher) publishWithRetry(ctx context.Context, msg Message) error {
	var err error
	for retry := 0; retry <= p.config.MaxRetries; retry++ {
		if err = p.executePublish(ctx, msg); err == nil {
			return nil
		}
		if retry < p.config.MaxRetries {
			time.Sleep(time.Second * time.Duration(retry+1))
		}
	}
	return err
}

func (p *Publisher) executePublish(ctx context.Context, msg Message) error {
	_, err := p.executor.ExecContext(ctx,
		"INSERT INTO outbox (uuid, event_topic, event_domain, event_type, object_type, producer, correlation_id, payload, metadata, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		msg.GetEventID(), msg.GetEventTopic(), msg.GetEventDomain(), msg.GetEventType(), msg.GetObjectType(),
		msg.GetProducer(), msg.GetCorrelationID(), msg.GetPayload(), msg.GetMetadata(), msg.GetCreatedAt())
	return err
}
