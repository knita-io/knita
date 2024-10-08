package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var defaultConfigFilePath = ""

func init() {
	home, _ := os.UserHomeDir()
	defaultConfigFilePath = filepath.Join(home, ".knita.yaml")
}

type config struct {
	Executors executorsConfig `mapstructure:"executors"`
}

type executorsConfig struct {
	Local  localExecutorConfig    `mapstructure:"local"`
	Remote []remoteExecutorConfig `mapstructure:"remote"`
}

type localExecutorConfig struct {
	// Disabled determines if the Knita CLI will run local builds.
	// If true, a local or upstream broker must be configured.
	Disabled bool `mapstructure:"disabled"`
	// Labels the executor will advertise to the broker.
	Labels []string `mapstructure:"labels"`
}

type remoteExecutorConfig struct {
	Disabled bool   `mapstructure:"disabled"`
	Address  string `mapstructure:"address"`
}

func fillDefaultValues(config *config) *config {
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
