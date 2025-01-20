package client

type Config struct {
	Timeout                   int  `yaml:"timeout"`
	MaxIdleConnections        int  `yaml:"max_idle_connections"`
	MaxConnectionsPerHost     int  `yaml:"max_connections_per_host"`
	MaxIdleConnectionsPerHost int  `yaml:"max_idle_connections_per_host"`
	AllowFollowRedirect       bool `yaml:"allow_follow_redirect"`
}
