package main

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type config struct {
	Broker localBrokerConfig `mapstructure:"broker"`
}

type localBrokerConfig struct {
	DisableLocalExecutor bool                  `mapstructure:"disable_local_executor"`
	Executors            []*executorConfig     `mapstructure:"executors"`
	Upstream             *upstreamBrokerConfig `mapstructure:"upstream"`
}

type executorConfig struct {
	Name     string `mapstructure:"name"`
	Address  string `mapstructure:"address"`
	Disabled bool   `mapstructure:"disabled"`
}

type upstreamBrokerConfig struct {
	Address string `mapstructure:"address"`
}

func getConfig(syslog *zap.SugaredLogger) (*config, error) {
	v := viper.New()
	v.SetConfigName(".knita")
	v.SetConfigType("yaml")
	v.AddConfigPath("$HOME")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	var notFoundErr viper.ConfigFileNotFoundError
	if err != nil {
		if errors.As(err, &notFoundErr) {
			return &config{}, nil
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	syslog.Infof("using config file: %s", v.ConfigFileUsed())
	conf := &config{}
	err = v.Unmarshal(conf)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config file: %w", err)
	}
	return conf, nil
}
