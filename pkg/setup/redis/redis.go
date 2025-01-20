package redis

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	setupLoggingKey   = "setup"
	setupLoggingValue = "redis"
)

// Redis - Redis storage.
type Redis struct {
	client *redis.Client
	log    *logrus.Logger
	stop   <-chan os.Signal
}

// NewRedis - Returns *Redis.
func NewRedis(cfg *Config, stop <-chan os.Signal, log *logrus.Logger) *Redis {
	client := getRedisClient(cfg)

	log.WithField(setupLoggingKey, setupLoggingValue).
		Info("NewRedis - redis client has been initialized")

	r := &Redis{
		client: client,
		log:    log,
	}

	go func(){
		if err := r.shutdown(); err != nil {
			log.WithFields(logrus.Fields{
				setupLoggingKey: setupLoggingValue,
			}).Error("r.shutdown error: ", err)
		}
	}()

	return r
}

// getRedisClient - Returns *redis.Client.
func getRedisClient(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv(cfg.URLEnvKey),
		Password: os.Getenv(cfg.PasswordEnvKey),
		DB: func() int {
			db, err := strconv.Atoi(os.Getenv(cfg.DBEnvKey))
			if err != nil {
				panic(err)
			}
			return db
		}(),
	})
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}

func (r *Redis) shutdown() error {
	select {
	case <-r.stop:
		r.log.WithField(setupLoggingKey, setupLoggingValue).
		Info("Shutdown - close redis connect")

		return r.client.Close()
	}
}
