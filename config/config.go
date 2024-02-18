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
}

type Web struct {
	ServerPort int `yaml:"serverPort"`
}

func InitConfig() {
	workdir, _ := os.Getwd()
	fmt.Println("work dir is:", workdir)
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
	mConfig := GlobalConfig.DB["default"]
	fmt.Println(mConfig.InvertIndexName, mConfig.IndexStorageDir, mConfig.ShardNum)
}
