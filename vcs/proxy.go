package vcs

import (
	"strings"

	"github.com/4396/pkg/vcs/parser"
	"github.com/4396/pkg/vcs/repository"
)

type pkg struct {
	Repo string
	Base string
}

var (
	pkgs = make(map[string]pkg)
)

func Proxy(name, repo, base string) {
	pkgs[name] = pkg{repo, base}
}

func proxyParse(path string) (string, string, error) {
	return githubParse(path)
}

func proxyRepository(pkg string) (repo, base string, err error) {
	for k, v := range pkgs {
		if strings.HasPrefix(pkg, k) {
			repo = v.Repo
			base = v.Base
			return
		}
	}

	err = ErrInvalidPackage
	return
}

func init() {
	match := func(s string) bool {
		for k := range pkgs {
			if strings.HasPrefix(s, k) {
				return true
			}
		}
		return false
	}
	parser.Registerf(match, proxyParse, true)
	repository.Registerf(match, proxyRepository, true)
}
