package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Authorization map[string]struct {
		ApiKey string
		Secret string
	}
}

// TODO watch config file
func Config() *Configuration {
	config := new(Configuration)
	configFile, err := os.Open("config.json")
	if err == nil {
		decoder := json.NewDecoder(configFile)
		if err := decoder.Decode(config); err == nil {
			return config
		}
	}
	return nil
}
