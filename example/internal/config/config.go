package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/woyow/setupper/pkg/setup/app"
	"github.com/woyow/setupper/pkg/setup/echotron"
	"github.com/woyow/setupper/pkg/setup/logger"
	"github.com/woyow/setupper/pkg/setup/psql"

	goEnv "github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	configDir = "configs"
)

// Config - Aggregate configurations for application.
type Config struct {
	App    app.Config    `yaml:"app"`
	Logger logger.Config `yaml:"logger"`
	Psql   psql.Config   `yaml:"psql"`
	YourTg struct {
		Echotron echotron.Config `yaml:"echotron"`
	} `yaml:"your_tg"`
}

// NewConfig - Returns *Config.
func NewConfig() (*Config, error) {
	var cfg Config

	if err := cfg.readEnv(); err != nil {
		return nil, err
	}

	if err := cfg.readConfigFile(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) readEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("config: readEnv - godotenv.Load error: ", err)
		return err
	}

	if _, err := goEnv.UnmarshalFromEnviron(cfg); err != nil {
		log.Fatal("config: readEnv - goEnv.UnmarshalFromEnviron error: ", err)
		return err
	}

	return nil
}

func (cfg *Config) readConfigFile() error {
	log.Println("config: readConfigFile - ", cfg.App.Env)
	fileName := configDir + "/" + cfg.App.Env + ".yaml"

	filePath, err := filepath.Abs(fileName)
	if err != nil {
		log.Fatal("config: readConfigFile - filepath.Abs error: ", err)
		return err
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("config: readConfigFile - os.ReadFile error: ", err)
		return err
	}

	if err = yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatal("config: readConfigFile - yaml.Unmarshal error: ", err)
		return err
	}

	return err
}
