package logger

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
	"net/http"
	"os"
	"time"
)

const (
	levelLoggingKey = "level"
)

func NewLogger(cfg *Config) *logrus.Logger {
	logger := logrus.StandardLogger()

	formatter := logrus.TextFormatter{
		DisableTimestamp: cfg.DisableTimestamp,
		FullTimestamp:    cfg.FullTimestamp,
	}

	logger.SetFormatter(&formatter)

	// Possible logLevel value: "panic", "fatal", "error", "warn" or "warning", "info", "debug", "trace"
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logger.WithError(err).
			WithField(levelLoggingKey, cfg.Level).
			Warn("cannot parse a logging level")
	} else {
		logger.SetLevel(level)
	}

	if cfg.Elastic.Enable {
		caCert, err := os.ReadFile(cfg.Elastic.Cert)
		if err != nil {
			logger.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		httpClient := &http.Client{
			Timeout: time.Duration(10) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		}

		client, err := elastic.NewClient(
			elastic.SetHttpClient(httpClient),
			elastic.SetURL(cfg.Elastic.URL),
			elastic.SetBasicAuth(
				os.Getenv(cfg.Elastic.UsernameEnvKey), 
				os.Getenv(cfg.Elastic.PasswordEnvKey),
			),
			elastic.SetHealthcheck(true),
			elastic.SetSniff(false),
		)
		if err != nil {
			logger.Panic(err)
		}

		hook, err := elogrus.NewAsyncElasticHook(client, "localhost", level, cfg.Elastic.IndexName)
		if err != nil {
			logger.Panic(err)
		}
		logger.Hooks.Add(hook)
	}

	return logger
}
