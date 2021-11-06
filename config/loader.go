package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// GlobalConfig ...
var GlobalConfig Config

// LoadConfig reads a YAML file and parses it into the GlobalConfig variable.
func LoadConfig(path string) {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal().Msgf("Failed to read yaml file", err)
	}

	err = yaml.Unmarshal(file, &GlobalConfig)
	if err != nil {
		log.Fatal().Msgf("Failed to parse yaml", err)
	}
}
