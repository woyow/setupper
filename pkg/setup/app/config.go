package app

type Config struct {
	Env     string `env:"APP_ENV,default=local"`
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
