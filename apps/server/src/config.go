package main

import (
	"peekaping/src/config"
)

func ProvideConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig("../..")

	if err != nil {
		panic(err)
	}

	return &cfg, nil
}
