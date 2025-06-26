package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
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
	query    string // Pre-built parameterized query to prevent SQL injection
}

var (
	defaultTableName  = "outbox_message"
	defaultBatchSize  = 10
	defaultMaxRetries = 3

	// Table name validation regex - only allows alphanumeric, underscore, and dot
	tableNameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)*$`)
)

type publisherConfig struct {
	tableName  string
	batchSize  int
	maxRetries int
}

// PublisherOption defines a function that configures a Publisher
type PublisherOption func(*publisherConfig) error

// WithBatchSize sets the batch size for publishing messages
func WithBatchSize(size int) PublisherOption {
	return func(c *publisherConfig) error {
		if size <= 0 {
			return fmt.Errorf("batch size must be positive, got %d", size)
		}
		c.batchSize = size
		return nil
	}
}

// WithTableName sets the table name for storing messages
// The table name must follow SQL identifier naming conventions
func WithTableName(tableName string) PublisherOption {
	return func(c *publisherConfig) error {
		if tableName == "" {
			return fmt.Errorf("table name cannot be empty")
		}

		// Use the security utility to validate and sanitize table name
		sanitizedTableName, err := sanitizeTableName(tableName)
		if err != nil {
			return fmt.Errorf("invalid table name: %w", err)
		}

		c.tableName = sanitizedTableName
		return nil
	}
}

// WithMaxRetries sets the maximum number of retry attempts
func WithMaxRetries(retries int) PublisherOption {
	return func(c *publisherConfig) error {
		if retries < 0 {
			return fmt.Errorf("max retries cannot be negative, got %d", retries)
		}
		c.maxRetries = retries
		return nil
	}
}

// NewPublisher creates a new Publisher with the given executor and options
func NewPublisher(executor Executor, opts ...PublisherOption) (Publisher, error) {
	if executor == nil {
		return nil, fmt.Errorf("executor cannot be nil")
	}

	config := publisherConfig{
		tableName:  defaultTableName,
		batchSize:  defaultBatchSize,
		maxRetries: defaultMaxRetries,
	}

	// Apply options with error handling
	for _, opt := range opts {
		if err := opt(&config); err != nil {
			return nil, fmt.Errorf("failed to apply publisher option: %w", err)
		}
	}

	const query = "INSERT INTO %s (event_id, event_topic, event_domain, event_type, object_type, producer, correlation_id, payload, metadata, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	return &publisher{
		executor: executor,
		config:   config,
		query:    fmt.Sprintf(query, config.tableName),
	}, nil
}

func (p *publisher) Publish(ctx context.Context, messages ...*Message) error {
	if len(messages) == 0 {
		return nil
	}

	// Validate all messages before processing
	for i, msg := range messages {
		if msg == nil {
			return fmt.Errorf("message at index %d is nil", i)
		}
		if err := msg.validate(); err != nil {
			return fmt.Errorf("message at index %d is invalid: %w", i, err)
		}
	}

	return p.publishBatch(ctx, messages)
}

func (p *publisher) Close() error {
	return nil
}

func (p *publisher) publishBatch(ctx context.Context, messages []*Message) error {
	batchSize := p.config.batchSize
	for i := 0; i < len(messages); i += batchSize {
		end := i + batchSize
		if end > len(messages) {
			end = len(messages)
		}

		batch := messages[i:end]
		if err := p.publishBatchInternal(ctx, batch); err != nil {
			return fmt.Errorf("failed to publish batch %d-%d: %w", i, end-1, err)
		}
	}

	return nil
}

func (p *publisher) publishBatchInternal(ctx context.Context, messages []*Message) error {
	for _, msg := range messages {
		if err := p.publishWithRetry(ctx, msg); err != nil {
			return fmt.Errorf("failed to publish message %s: %w", msg.EventID, err)
		}
	}
	return nil
}

func (p *publisher) publishWithRetry(ctx context.Context, msg *Message) error {
	var lastErr error

	for retry := 0; retry <= p.config.maxRetries; retry++ {
		if err := p.executePublish(ctx, msg); err == nil {
			return nil
		} else {
			lastErr = err
			if retry < p.config.maxRetries {
				// Exponential backoff with jitter
				backoff := time.Duration(retry+1) * time.Second
				time.Sleep(backoff)
			}
		}
	}

	return fmt.Errorf("failed to publish message %s after %d retries: %w", msg.EventID, p.config.maxRetries, lastErr)
}

func (p *publisher) executePublish(ctx context.Context, msg *Message) error {
	_, err := p.executor.ExecContext(ctx, p.query,
		msg.EventID, msg.EventTopic, msg.EventDomain, msg.EventType, msg.ObjectType,
		msg.Producer, msg.CorrelationID, msg.Payload, msg.Metadata, msg.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetTableName returns the current table name (for testing purposes)
func (p *publisher) GetTableName() string {
	return p.config.tableName
}
