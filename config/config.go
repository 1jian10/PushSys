package config

import "github.com/spf13/viper"

type Config struct {
	Websocket struct {
		WriteBufferSize int
		ReadBufferSize  int
		Port            string
	}
	Etcd struct {
		EndPoints   []string
		DialTimeout int64
		Name        string
		Addr        string
		TTL         int64
	}
	Redis struct {
		Addr string
		DB   int
	}
}

func ReadConfig() Config {
	var c Config
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	if err := viper.Unmarshal(&c); err != nil {
		panic(err.Error())
	}
	return c
}
