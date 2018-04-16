package vcs

import (
	"errors"
	"strings"

	"github.com/4396/pkg/vcs/archive"
	"github.com/4396/pkg/vcs/parser"
	"github.com/4396/pkg/vcs/repository"
	"github.com/4396/pkg/vcs/revision"
)

var (
	ErrNoService      = errors.New("no service")
	ErrInvalidSHA     = errors.New("invalid sha")
	ErrInvalidPath    = errors.New("invalid path")
	ErrInvalidRepo    = errors.New("invalid repo")
	ErrInvalidPackage = errors.New("invalid package")
)

func Archive(repo, sha string) (url string, err error) {
	srv, ok := archive.Select(repo)
	if !ok {
		err = ErrNoService
		return
	}

	url, err = srv.Archive(repo, sha)
	return
}

func Repository(pkg string) (repo, base string, err error) {
	srv, ok := repository.Select(pkg)
	if !ok {
		err = ErrNoService
		return
	}

	repo, base, err = srv.Repository(pkg)
	return
}

func Revision(repo, ver string) (sha string, err error) {
	srv, ok := revision.Select(repo)
	if !ok {
		err = ErrNoService
		return
	}

	sha, err = srv.Revision(repo, ver)
	return
}

func Parse(path string) (pkg, ver string, err error) {
	srv, ok := parser.Select(path)
	if !ok {
		err = ErrNoService
		return
	}

	pkg, ver, err = srv.Parse(path)
	return
}

func prefixMatchFunc(name string) func(string) bool {
	return func(s string) bool {
		return strings.HasPrefix(s, name)
	}
}
