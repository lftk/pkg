package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/4396/pkg/store"
	"gopkg.in/yaml.v2"
)

var (
	err  error
	once sync.Once
	pkg  = config{
		Registry: defaultRegistry,
		Token:    defaultToken,
	}
	defaultRegistry = "https://pkg.4396.io"
	defaultToken    = "0443dbd565c01d39cb97a4e452d580986251d6c5"
)

type config struct {
	Registry string `yaml:"registry"`
	Token    string `yaml:"token"`
}

func loadConfig() error {
	once.Do(func() {
		var b []byte
		dir, _ := store.Dir("")
		path := filepath.Join(dir, ".config")
		_, errStat := os.Stat(path)
		if errStat != nil {
			return
		}

		b, err = ioutil.ReadFile(path)
		if err != nil {
			return
		}

		err = yaml.Unmarshal(b, &pkg)
	})
	return err
}

func saveConfig() (err error) {
	b, err := yaml.Marshal(&pkg)
	if err != nil {
		return
	}

	dir, _ := store.Dir("")
	path := filepath.Join(dir, ".config")
	err = ioutil.WriteFile(path, b, os.ModePerm)
	return
}

func Registry(url string) (oldURL string, err error) {
	err = loadConfig()
	if err != nil {
		return
	}

	oldURL = pkg.Registry
	if url == "" || url == pkg.Registry {
		return
	}

	pkg.Registry = url
	err = saveConfig()
	return
}

func Token(token string) (oldToken string, err error) {
	err = loadConfig()
	if err != nil {
		return
	}

	oldToken = pkg.Token
	if token == "" || token == pkg.Token {
		return
	}

	pkg.Token = token
	err = saveConfig()
	return
}
