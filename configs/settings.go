package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	ListenAddress    string
	ListenPort       string
	TopDeviceRole    string
	BottomDeviceRole string
	Backend          string
	Backends         struct {
		File struct {
			Path string
		}
		Redis struct {
			Host     string
			Port     int
			Password string
			DB       int
		}
	}
}

var Config config

func LoadConfig() error {
	viper.SetConfigName("settings")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("claw")
	viper.AutomaticEnv()

	viper.SetDefault("ListenAddress", "0.0.0.0")
	viper.SetDefault("ListenPort", "8080")
	viper.SetDefault("TopDeviceRole", "edge")
	viper.SetDefault("BottomDeviceRole", "tor")

	viper.SetDefault("Backend", "File")

	viper.SetDefault("Backends.File.Path", "./topologies/")

	viper.SetDefault("Backends.Redis.Host", "localhost")
	viper.SetDefault("Backends.Redis.Port", 6379)
	viper.SetDefault("Backends.Redis.Password", "")
	viper.SetDefault("Backends.Redis.DB", 0)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return err
	}

	return nil
}
