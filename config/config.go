package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Db   string
	User UserConfig
}

type UserConfig struct {
	FirstName string
	LastName  string
	Fonction  string
}

const (
	ConfigFile     = ".dbconfig"
	DefaultDbAlias = "default"
	UserDbAlias    = "user"
)

var NoConfig = errors.New(fmt.Sprintf("no %s found", ConfigFile))

// Find the closest directory containing .dbconfig starting
// in cwd and then searching up ancestor tree.
// Look first looking in cwd and then up through its ancestors
func FindConfig() (*Config, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for {
		conf := filepath.Join(curDir, ConfigFile)
		info, err := os.Stat(conf)
		if err == nil && !info.IsDir() {
			// found
			return ReadConfig(conf)
		} else if err != nil && !os.IsNotExist(err) {
			// can't read
			return nil, err
		}
		nextDir := filepath.Dir(curDir)
		if nextDir == curDir {
			// stop at root
			return nil, NoConfig
		}
		curDir = nextDir
	}
}

func ReadConfig(name string) (*Config, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var conf Config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		return nil, err
	}
	return &conf, err
}
