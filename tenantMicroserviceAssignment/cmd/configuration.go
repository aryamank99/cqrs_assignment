package cmd

// These variables are meant to inserted as env vars
// you can inject them through orchestrator secrets or command line args or OS env vars
// Precedence: Command-line flags > Env vars > Default values.

type Config struct {
	Environment string `arg:"env:ENVIRONMENT"`
	ServerConfig
	CommandDatabaseConfig
	StoreDatabaseConfig
	AMQPConfig
}

type ServerConfig struct {
	Port string `arg:"env:SERVER_PORT"`
	Name string `arg:"env:SERVER_NAME"`
}

type AMQPConfig struct {
	RabbitMQConnectionString string `arg:"env:RABBIT_MQ_CONNECTION_STRING"`
}

type StoreDatabaseConfig struct {
	DatabaseName           string `arg:"env:DATABASE_NAME"`
	DatabaseUsername       string `arg:"env:DATABASE_USERNAME"`
	DatabasePassword       string `arg:"env:DATABASE_PASSWORD"`
	DatabaseConnectionUrl  string `arg:"env:DATABASE_CONNECTION_URL"`
	DatabaseConnectionPort string `arg:"env:DATABASE_CONNECTION_PORT"`
}

type CommandDatabaseConfig struct {
	DatabaseName           string `arg:"env:DATABASE_NAME"`
	DatabaseUsername       string `arg:"env:DATABASE_USERNAME"`
	DatabasePassword       string `arg:"env:DATABASE_PASSWORD"`
	DatabaseConnectionUrl  string `arg:"env:DATABASE_CONNECTION_URL"`
	DatabaseConnectionPort string `arg:"env:DATABASE_CONNECTION_PORT"`
}

// DefaultConfiguration return dev configuration as default config
func DefaultConfiguration() *Config {
	return &Config{
		Environment: "dev",
		ServerConfig: ServerConfig{
			Name: "go-lang-server",
			Port: "8090",
		},
		StoreDatabaseConfig: StoreDatabaseConfig{
			DatabaseName:           "tenant_db_store",
			DatabaseUsername:       "tenant_user",
			DatabasePassword:       "pass1234",
			DatabaseConnectionUrl:  "localhost",
			DatabaseConnectionPort: "3306",
		},
		CommandDatabaseConfig: CommandDatabaseConfig{
			DatabaseName:           "tenant_db_commands",
			DatabaseUsername:       "tenant_user",
			DatabasePassword:       "pass1234",
			DatabaseConnectionUrl:  "localhost",
			DatabaseConnectionPort: "3306",
		},
		AMQPConfig: AMQPConfig{
			RabbitMQConnectionString: "amqp://guest:guest@localhost:5672/",
		},
	}
}
