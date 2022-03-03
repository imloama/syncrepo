package sync

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Ip string
	Port int
	Folder string
	Repos []Repo
}

type Repo struct {
	Name     string // 名称
	From     string // 来源仓库
	Target   string // 目的仓库
	Branch   string // 分支名
	Interval int    // 时间间隔，单位是秒
	Enable   bool   // 是否启用
}

var config Config
var repos = make(map[string]*Repo)

func initConfig(){
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")               // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	if err = viper.Unmarshal(&config); err!=nil{
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	if len(config.Repos) <= 0 {
		return
	}
	for _, repo := range config.Repos {
		repos[repo.Name] = &repo
	}
}

func GetConfig()*Config{
	return &config
}

func DumpConfig(){
	data, err := json.Marshal(config)
	if err!=nil{
		panic(fmt.Errorf("Fatal error config: %w \n", err))
	}
	fmt.Printf("%v \n", string(data))
}