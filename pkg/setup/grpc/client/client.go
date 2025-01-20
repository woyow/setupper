package client

import (
	"os"

	"github.com/processout/grpc-go-pool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pool *grpcpool.Pool
}

func NewClient(cfg *Config, log *logrus.Logger) (*Client, error) {
	pool, err := grpcpool.New(func() (*grpc.ClientConn, error) {
		address := os.Getenv(cfg.AddressEnvKey)
		if address == "" {
			panic(ErrEmptyAddress)
		}
		return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}, 10, 10, 10, 0)
	if err != nil {
		log.Error("setup grpc client: NewClient - grpcpool.New error: ", err.Error())
		return nil, err
	}

	return &Client{
		pool: pool,
	}, nil
}

func (c *Client) Pool() *grpcpool.Pool {
	return c.pool
}
