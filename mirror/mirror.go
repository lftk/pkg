package mirror

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/4396/pkg/store"
	"github.com/4396/pkg/vcs"
	"github.com/ghodss/yaml"
)

type mirror struct {
	Repository string `yaml:"repository"`
	Base       string `yaml:"base"`
}

var (
	mirrors = make(map[string]mirror)
)

func init() {
	dir, _ := store.Dir("")
	path := filepath.Join(dir, ".mirror")
	_, err := os.Stat(path)
	if err != nil {
		return
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(b, &mirrors)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range mirrors {
		vcs.Proxy(k, v.Repository, v.Base)
	}
}

func saveMirrors() (err error) {
	b, err := yaml.Marshal(&mirrors)
	if err != nil {
		return
	}

	dir, _ := store.Dir("")
	path := filepath.Join(dir, ".mirror")
	err = ioutil.WriteFile(path, b, os.ModePerm)
	return
}

func Set(pkg, repo, base string) (err error) {
	if v, ok := mirrors[pkg]; ok {
		if v.Repository == repo && v.Base == base {
			return
		}
	}

	mirrors[pkg] = mirror{
		Repository: repo,
		Base:       base,
	}
	err = saveMirrors()
	return
}

func Get(pkg string) (repo, base string, ok bool) {
	v, ok := mirrors[pkg]
	if ok {
		repo = v.Repository
		base = v.Base
	}
	return
}

func Delete(pkg string) {
	if _, ok := mirrors[pkg]; ok {
		delete(mirrors, pkg)
		saveMirrors()
	}
	return
}
