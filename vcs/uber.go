package vcs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/4396/pkg/vcs/parser"
	"github.com/4396/pkg/vcs/repository"
)

func uberParse(path string) (string, string, error) {
	return githubParse(path)
}

func uberRepository(pkg string) (repo, base string, err error) {
	ss := strings.Split(pkg, "/")
	if len(ss) < 2 {
		err = ErrInvalidPackage
		return
	}

	repo = fmt.Sprintf("github.com/uber-go/%s", ss[1])
	base = filepath.Join("go.uber.org", ss[1])
	return
}

func init() {
	match := prefixMatchFunc("go.uber.org")
	repository.Registerf(match, uberRepository, false)
	parser.Registerf(match, uberParse, false)
}
