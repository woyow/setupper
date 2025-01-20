package psql

// Config - Postgresql config
type Config struct {
	URLEnvKey           string    `yaml:"url_env_key" default:"PG_URL"`
	DBNameEnvKey        string    `yaml:"db_name_env_key" default:"PG_DATABASE"`
	URLParametersEnvKey string    `yaml:"url_parameters_env_key" default:"PG_URL_PARAMETERS"`
	Migration           Migration `yaml:"migration"`
	Pool                Pool      `yaml:"pool"`
	Host                string    `yaml:"host"`
	Port                string    `yaml:"port"`
	SSLMode             string    `yaml:"sslmode"`
}


type Migration struct {
	Source          string `yaml:"source"`
	Attempts        int    `yaml:"attempts"`
	AttemptsTimeout int    `yaml:"attempts_timeout"`
	Enable          bool   `yaml:"enable"`
}

// Pool - Pool config
type Pool struct {
	MaxPoolSize        int `yaml:"max_pool_size"`
	MinPoolSize        int `yaml:"min_pool_size"`
	ConnectionAttempts int `yaml:"connection_attempts"`
	ConnectionTimeout  int `yaml:"connection_timeout"`
}
