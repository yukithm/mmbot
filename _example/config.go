package main

import (
	"io/ioutil"

	"github.com/naoina/toml"
	"github.com/yukithm/mmbot/app"
)

type appConfig struct {
	app.Config
	Example struct {
		Foo int    `toml:"foo"`
		Bar string `toml:"bar"`
	} `toml:"example"`
}

func loadConfig(file string) (*appConfig, error) {
	var config appConfig
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
