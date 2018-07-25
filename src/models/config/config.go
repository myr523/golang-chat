package config

import (
	"encoding/json"
	"io/ioutil"
)

type GoogleConfig struct {
	CLIENTID     string
	CLIENTSECRET string
}

type GomniConfig struct {
	SECURITYKEY string
}

type ConfigAll struct {
	GOOGLE    GoogleConfig
	GOMNIAUTH GomniConfig
}

func Perse(filename string) (*ConfigAll, error) {
	var config ConfigAll
	jsonStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonStr, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
