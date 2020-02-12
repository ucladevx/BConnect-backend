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
		Port     string `yaml:"PORT"`
		Host     string `yaml:"HOST"`
		Username string `yaml:"USERNAME"`
		Name     string `yaml:"NAME"`
		Password string `yaml:"PASSWORD"`
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
