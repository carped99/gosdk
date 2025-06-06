package outbox

import (
	"context"
	"database/sql"
	"time"
)

type Publisher interface {
	Publish(ctx context.Context, messages ...*Message) error
	Close() error
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type publisher struct {
	executor Executor
	config   publisherConfig
}

var (
	defaultTableName  = "outbox_message"
	defaultBatchSize  = 10
	defaultMaxRetries = 3
)

type publisherConfig struct {
	tableName  string
	batchSize  int
	maxRetries int
}

// PublisherOption defines a function that configures a Publisher
type PublisherOption func(*publisherConfig)

// WithBatchSize sets the batch size for publishing messages
func WithBatchSize(size int) PublisherOption {
	return func(c *publisherConfig) {
		if size <= 0 {
			c.batchSize = defaultBatchSize // default batch size
		}
		c.batchSize = size
	}
}

func WithTableName(tableName string) PublisherOption {
	return func(c *publisherConfig) {
		if tableName == "" {
			c.tableName = defaultTableName // default table name
		} else {
			c.tableName = tableName
		}
	}
}

// WithMaxRetries sets the maximum number of retry attempts
func WithMaxRetries(retries int) PublisherOption {
	return func(c *publisherConfig) {
		if retries <= 0 {
			c.maxRetries = defaultMaxRetries // default max retries
		}
		c.maxRetries = retries
	}
}

// NewPublisher creates a new Publisher with the given executor and options
func NewPublisher(executor Executor, opts ...PublisherOption) Publisher {
	config := publisherConfig{
		tableName:  defaultTableName,
		batchSize:  defaultBatchSize,  // default batch size
		maxRetries: defaultMaxRetries, // default max retries
	}

	// Apply options
	for _, opt := range opts {
		opt(&config)
	}

	return &publisher{
		executor: executor,
		config:   config,
	}
}

func (p *publisher) Publish(ctx context.Context, messages ...*Message) error {
	if len(messages) == 0 {
		return nil
	}
	return p.publishBatch(ctx, messages)
}

func (p *publisher) Close() error {
	// Implement any necessary cleanup logic here, if needed
	return nil
}

func (p *publisher) publishBatch(ctx context.Context, messages []*Message) error {
	for _, msg := range messages {
		if err := p.publishWithRetry(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func (p *publisher) publishWithRetry(ctx context.Context, msg *Message) error {
	var err error
	for retry := 0; retry <= p.config.maxRetries; retry++ {
		if err = p.executePublish(ctx, msg); err == nil {
			return nil
		}
		if retry < p.config.maxRetries {
			time.Sleep(time.Second * time.Duration(retry+1))
		}
	}
	return err
}

func (p *publisher) executePublish(ctx context.Context, msg *Message) error {
	_, err := p.executor.ExecContext(ctx,
		"INSERT INTO "+p.config.tableName+" (event_id, event_topic, event_domain, event_type, object_type, producer, correlation_id, payload, metadata, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		msg.EventID, msg.EventTopic, msg.EventDomain, msg.EventType, msg.ObjectType,
		msg.Producer, msg.CorrelationID, msg.Payload, msg.Metadata, msg.CreatedAt)
	return err
}
