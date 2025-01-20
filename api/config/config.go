package config

import (
	"github.com/spf13/viper"
	"strings"
)

func init() {
	viper.SetDefault("httpport", 80)
	viper.SetDefault("data.dir", "/Users/user/Data")
}

type Config struct {
	HttpPort int
	Data     Data
}

type Data struct {
	Dir string
}

func InitializeConfig() (*Config, error) {
	v := viper.GetViper()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
