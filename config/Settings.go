package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Api      Api      `json:"api"`
	Database Database `json:"database"`
	Auth     Auth     `json:"auth"`
}

func NewSettings(filename string) (*Settings, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	yamlDecoder := yaml.NewDecoder(file)

	settings := new(Settings)
	err = yamlDecoder.Decode(&settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
