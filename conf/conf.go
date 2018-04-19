package conf

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
	cfg  = config{
		Registry: DefaultRegistry,
		Token:    DefaultToken,
	}
	DefaultRegistry = "https://pkg.4396.io"
	DefaultToken    = "0443dbd565c01d39cb97a4e452d580986251d6c5"
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

		err = yaml.Unmarshal(b, &cfg)
	})
	return err
}

func saveConfig() (err error) {
	b, err := yaml.Marshal(&cfg)
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

	oldURL = cfg.Registry
	if url == "" || url == cfg.Registry {
		return
	}

	cfg.Registry = url
	err = saveConfig()
	return
}

func Token(token string) (oldToken string, err error) {
	err = loadConfig()
	if err != nil {
		return
	}

	oldToken = cfg.Token
	if token == "" || token == cfg.Token {
		return
	}

	cfg.Token = token
	err = saveConfig()
	return
}
