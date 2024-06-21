package main

import (
	"fmt"
	"github.com/rs/xid"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type config struct {
	BindAddress string `mapstructure:"bind_address"`
	// Name is the human-friendly name of the executor.
	// This should be unique across all executors.
	Name string `mapstructure:"name"`
	// Labels the executor will advertise to the broker.
	Labels []string `mapstructure:"labels"`
}

func fillDefaultValues(config *config) *config {
	if config.BindAddress == "" {
		config.BindAddress = "127.0.0.1:9091"
	}
	if config.Name == "" {
		host, _ := os.Hostname()
		if host == "" {
			config.Name = fmt.Sprintf("%s (name unconfigured)", xid.New().String())
		} else {
			config.Name = host
		}
	}
	return config
}

func getConfig(syslog *zap.SugaredLogger, configFilePath string) (*config, error) {
	v := viper.New()
	v.AutomaticEnv()
	_, err := os.Stat(configFilePath)
	if err == nil {
		v.SetConfigFile(configFilePath)
		err := v.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		syslog.Infof("using config file: %s", v.ConfigFileUsed())
	}
	conf := &config{}
	err = v.Unmarshal(conf)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config file: %w", err)
	}
	return fillDefaultValues(conf), nil
}
