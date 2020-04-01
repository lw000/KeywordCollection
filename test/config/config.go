package config

import (
	"encoding/json"
	"io/ioutil"
)

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}

// Load ...
func (c *Config) Load(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}
