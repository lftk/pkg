package vcs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/4396/pkg/vcs/parser"
	"github.com/4396/pkg/vcs/repository"
)

func golangParse(path string) (string, string, error) {
	return githubParse(path)
}

func golangRepository(pkg string) (repo, base string, err error) {
	ss := strings.Split(pkg, "/")
	if len(ss) < 3 {
		err = ErrInvalidRepo
		return
	}

	repo = fmt.Sprintf("github.com/golang/%s", ss[2])
	base = filepath.Join(ss[:3]...)
	return
}

func init() {
	match := prefixMatchFunc("golang.org")
	parser.Registerf(match, golangParse, false)
	repository.Registerf(match, golangRepository, false)
}
