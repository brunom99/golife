package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server struct {
		Port int
	}
	Grid struct {
		Size int
		Seed int64
	}
	Bubble struct {
		Proba float64
		Pool  int
	}
	Bubbles map[string]struct {
		Pool     int
		MinSpeed int
		MaxSpeed int
		Diagonal bool
	}
}

func (c *Config) LoadFile(tomlFile string) error {
	// open toml file
	file, err := os.Open(tomlFile)
	if err != nil {
		return err
	}
	defer file.Close()
	// read toml file
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	// unmarshal toml
	if err = toml.Unmarshal(b, c); err != nil {
		return err
	}
	// done
	return nil
}
