package service

type Config struct {
	ServiceName string `env:"SERVICE_NAME"`

	DbGormDisabled bool `env:"DB_GORM_DISABLED"`

	SchedulerEnabled bool `env:"SCHEDULER_ENABLED"`

	RestServerHost string `env:"REST_SERVER_HOST" envDefault:"0.0.0.0"`
	RestServerPort int    `env:"REST_SERVER_PORT" envDefault:"8080"`

	GrpcServerDisabled    bool   `env:"GRPC_SERVER_DISABLED"`
	GrpcServerHost        string `env:"GRPC_SERVER_HOST" envDefault:"0.0.0.0"`
	GrpcServerPort        int    `env:"GRPC_SERVER_PORT" envDefault:"9000"`
	GrpcReflectionEnabled bool   `env:"GRPC_REFLECTION_ENABLED" envDefault:"true"`
}
