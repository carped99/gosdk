package entx

type DatabaseConfig struct {
	Driver          string `json:"driver" mapstructure:"driver"`
	DSN             string `json:"dsn" mapstructure:"dsn"`
	MaxOpenConns    int    `json:"maxOpenConns" mapstructure:"maxOpenConns"`
	MaxIdleConns    int    `json:"maxIdleConns" mapstructure:"maxIdleConns"`
	ConnMaxLifetime int    `json:"connMaxLifetime" mapstructure:"connMaxLifetime"`
}

var DefaultDatabaseConfig = DatabaseConfig{
	Driver: "postgres",
	DSN:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
}
