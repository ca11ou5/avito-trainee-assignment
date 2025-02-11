package configs

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	PostgresURL string `env:"POSTGRES_URL" env-required:"true"`

	PORT string `env:"PORT" env-required:"true"`
}
