package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	General   ConfigGeneral
	CIPClass3 []ConfigCIPClass3
}

func (c *Config) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("problem saving config to %s: %w", filename, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)

	return enc.Encode(*c)
}

func ConfigNew() *Config {
	c := new(Config)
	c.CIPClass3 = make([]ConfigCIPClass3, 0)
	return c
}

func ConfigLoad(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("problem opening config file: %w", err)
	}

	j := json.NewDecoder(f)
	c := new(Config)
	err = j.Decode(c)
	if err != nil {
		return nil, fmt.Errorf("problem parsing config file: %w", err)
	}

	return c, nil
}

type ConfigGeneral struct {
	Host string
	Port int
}
