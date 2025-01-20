package redis

type Config struct {
	URLEnvKey      string `yaml:"url_env_key" default:"REDIS_URL"`
	PasswordEnvKey string `yaml:"password_env_key" default:"REDIS_PASSWORD"`
	DBEnvKey       string `yaml:"db_env_key" default:"REDIS_DB"`
}
