package bootstrap

type BrokerType string

const (
	RabbitMQ BrokerType = "rabbitmq"
	Kafka    BrokerType = "kafka"
)

type BrokerConfig struct {
	Driver   BrokerType     `json:"type" mapstructure:"type"`
	Enabled  bool           `json:"enabled" mapstructure:"enabled"`
	Kafka    KafkaConfig    `json:"kafka" mapstructure:"kafka"`
	RabbitMQ RabbitMQConfig `json:"rabbitmq" mapstructure:"rabbitmq"`
	Memory   MemoryConfig   `json:"memory" mapstructure:"memory"`
}

var DefaultBrokerConfig = BrokerConfig{
	Driver:   RabbitMQ,
	Enabled:  true,
	Kafka:    KafkaConfig{},
	RabbitMQ: DefaultRabbitMQConfig,
	Memory:   MemoryConfig{},
}

type KafkaConfig struct {
	Brokers    []string  `json:"brokers"`
	Topic      string    `json:"topic"`
	TLSEnabled bool      `json:"tls_enabled"`
	SASL       KafkaSASL `json:"sasl"`
}

type KafkaSASL struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Mechanism string `json:"mechanism"`
}

type RabbitMQConfig struct {
	URI        string `json:"uri" mapstructure:"uri"`
	Exchange   string `json:"exchange" mapstructure:"exchange"`
	Queue      string `json:"queue" mapstructure:"queue"`
	RoutingKey string `json:"routing_key" mapstructure:"routing_key"`
}

var DefaultRabbitMQConfig = RabbitMQConfig{
	URI:        "amqp://admin:admin@localhost:5672",
	Exchange:   "mago",
	Queue:      "",
	RoutingKey: "",
}

type MemoryConfig struct {
	BufferSize int `json:"buffer_size"`
}
