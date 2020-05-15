package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/attic-labs/noms/go/spec"
)

type Config struct {
	Url  string
	User UserConfig
}

type UserConfig struct {
	FirstName string
	LastName  string
	Fonction  string
}

const (
	ConfigFile  = ".dbconfig"
	UserDbAlias = "user"
	ConfigDb    = "db"
)

type Configs struct {
	File string
	Conf map[string]Config
}

var NoConfig = errors.New(fmt.Sprintf("no %s found", ConfigFile))

// Find the closest directory containing .dbconfig starting
// in cwd and then searching up ancestor tree.
// Look first looking in cwd and then up through its ancestors
func FindConfig() (*Configs, error) {
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

func ReadConfig(name string) (*Configs, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	c, err := NewConfig(string(data))
	if err != nil {
		return nil, err
	}
	c.File = name
	return qualifyPaths(name, c)
}

func qualifyPaths(configPath string, c *Configs) (*Configs, error) {
	file, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(file)
	qc := *c
	qc.File = file
	for k, r := range c.Conf {
		qc.Conf[k] = Config{absDbSpec(dir, r.Url), r.User}
	}
	return &qc, nil
}

func NewConfig(data string) (*Configs, error) {
	c := new(Configs)
	if _, err := toml.Decode(string(data), &c.Conf); err != nil {
		return nil, err
	}
	return c, nil
}

// Replace relative directory in path part of spec with an absolute
// directory. Assumes the path is relative to the location of the config file
func absDbSpec(configHome string, url string) string {
	dbSpec, err := spec.ForDatabase(url)
	if err != nil {
		return url
	}
	if dbSpec.Protocol != "nbs" {
		return url
	}
	dbName := dbSpec.DatabaseName
	if !filepath.IsAbs(dbName) {
		dbName = filepath.Join(configHome, dbName)
	}
	return "nbs:" + dbName
}

func (c *Configs) WriteTo(configHome string) (string, error) {
	file := filepath.Join(configHome, ConfigFile)
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(file, []byte(c.writeableString()), os.ModePerm); err != nil {
		return "", err
	}
	return file, nil
}

func (c *Configs) writeableString() string {
	var buffer bytes.Buffer
	myConf := c.Conf["db"]
	buffer.WriteString(fmt.Sprintf("[db]\n"))
	buffer.WriteString(fmt.Sprintf(`Url = "%s"`+"\n", myConf.Url))
	buffer.WriteString(fmt.Sprintf("\t" + "[db.user]\n"))
	buffer.WriteString(fmt.Sprintf("\t"+`FirstName = "%s"`+"\n", myConf.User.FirstName))
	buffer.WriteString(fmt.Sprintf("\t"+`LastName = "%s"`+"\n", myConf.User.LastName))
	buffer.WriteString(fmt.Sprintf("\t"+`Fonction = "%s"`+"\n", myConf.User.Fonction))
	return buffer.String()
}
