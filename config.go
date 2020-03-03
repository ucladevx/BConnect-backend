package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

//Config configuration
type Config struct {
	Server struct {
		Port string `yaml:"PORT"`
		Host string `yaml:"HOST"`
	} `yaml:"SERVER"`
	Storage struct {
		Port           string `yaml:"PORT"`
		UserHost       string `yaml:"USERHOST"`
		FriendHost     string `yaml:"FRIENDHOST"`
		Friendname     string `yaml:"FRIENDNAME"`
		Username       string `yaml:"USERNAME"`
		UserUsername   string `yaml:"USERUSERNAME"`
		UserPassword   string `yaml:"USERPASSWORD"`
		FriendUsername string `yaml:"FRIENDUSERNAME"`
		FriendPassword string `yaml:"FRIENDPASSWORD"`
		HerokuUser     string `yaml:"HEROKUUSER"`
		HerokuFriend   string `yaml:"HEROKUFRIEND"`
	} `yaml:"STORAGE"`
}

//Conf config
func Conf() Config {
	f, err := os.Open("./credentials/config.yml")
	if err != nil {
		print(err.Error())
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		print(err.Error())
	}

	return cfg
}
