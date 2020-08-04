package utils

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var Config = LoadConfig()

type ConfigFile struct {
	UnbotKeys []string `yaml:"unbot_keys"`
	OwmKey    string   `yaml:"openweathermap_key"`
	Location  struct {
		Country   string
		Latitude  float64
		Longitude float64
	}
}

func LoadConfig() ConfigFile {
	var config ConfigFile

	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if yaml.Unmarshal(data, &config) != nil {
		log.Fatal("Error parsing config file")
	}

	return config
}
