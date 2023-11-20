package config

import (
	"context"

	envconfig "github.com/sethvargo/go-envconfig"
)

const (
	mongoDBPrefix = "MONGO_DB_"
)

type MongoDB struct {
	URI      string `env:"URI"`
	Password string `env:"PASSWORD"`
	User     string `env:"USER"`
	Database string `env:"DATABASE"`
}

func NewMongoDB(ctx context.Context) (*MongoDB, error) {
	conf := &MongoDB{}
	pl := envconfig.PrefixLookuper(mongoDBPrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		return nil, err
	}

	return conf, nil
}
