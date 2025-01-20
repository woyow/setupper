package logger

// Config - Logger config.
type Config struct {
	Elastic          Elastic `yaml:"elastic"`
	Level            string  `yaml:"log_level"`
	DisableTimestamp bool    `yaml:"disable_timestamp"`
	FullTimestamp    bool    `yaml:"full_timestamp"`
}

type Elastic struct {
	URL            string `yaml:"url"`
	IndexName      string `yaml:"index_name"`
	Cert           string `yaml:"cert"`
	UsernameEnvKey string `yaml:"username_env_key" default:"ELASTIC_USERNAME"`
	PasswordEnvKey string `yaml:"password_env_key" default:"ELASTIC_PASSWORD"`
	Enable         bool   `yaml:"enable"`
}
