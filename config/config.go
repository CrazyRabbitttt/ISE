package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var GlobalConfig *Config

type Config struct {
	DB  map[string]*DB  `yaml:"db"`
	Web map[string]*Web `yaml:"web"`
}

type DB struct {
	IndexStorageDir       string `yaml:"indexStorageDir"`
	BufferNum             int    `yaml:"bufferNum"`
	ShardNum              int    `yaml:"shardNum"`
	InvertIndexName       string `yaml:"invertIndexName"`
	PositiveIndexName     string `yaml:"positiveIndexName"`
	RepositoryStorageName string `yaml:"repositoryStorageName"`
	TimeOut               int64  `yaml:"timeOut"`
}

type Web struct {
	ServerPort int `yaml:"serverPort"`
}

func InitConfig() {
	workdir, _ := os.Getwd()
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workdir + "/config")
	viper.AddConfigPath(workdir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("Config 初始化成功")
}
