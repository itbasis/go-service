package db

type Config struct {
	Dialect string `env:"DB_DIALECT" envDefault:"pgx"`
	Host    string `env:"DB_HOST,notEmpty"`
	Port    int    `env:"DB_PORT" envDefault:"5432"`
	Name    string `env:"DB_NAME,notEmpty"`
	SslMode string `env:"DB_SSL_MODE" envDefault:"disable"`

	MaxIdleConnections   int `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"5"`
	MaxOpenConnections   int `env:"DB_MAX_OPEN_CONNECTIONS" envDefault:"15"`
	MaxLifetimeInMinutes int `env:"DB_MAX_LIFETIME_IN_MINUTES" envDefault:"15"`

	GooseMigrationDir string `env:"GOOSE_MIGRATION_DIR" envDefault:"/db/migrations"`
}

type Credential struct {
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
}
