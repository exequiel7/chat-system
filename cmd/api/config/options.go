package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Option func(*Config)

func WithConfigFile(path string, cname string, ctype string) Option {
	return func(c *Config) {
		fmt.Printf("%s/%s.%s", path, cname, ctype)
		viper.AddConfigPath(path)
		viper.SetConfigName(cname)
		viper.SetConfigType(ctype)
	}
}
