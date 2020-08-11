package pkg

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var Config = LoadConfig()

type ConfigFile struct {
	UnbotKeys []string `yaml:"keys"`
	OwmKey    string   `yaml:"openweathermap_key"`
	Country   string
	Places    struct {
		Default Place
		Names   []string
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
