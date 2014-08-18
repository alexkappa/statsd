package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Addr string `json:"address"`
}

func New(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}

func Decode(f *os.File) (*Config, error) {
	c := &Config{}
	err := json.NewDecoder(f).Decode(c)
	return c, err
}
