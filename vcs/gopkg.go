package vcs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/4396/pkg/vcs/parser"
	"github.com/4396/pkg/vcs/repository"
)

func gopkgParse(path string) (pkg, ver string, err error) {
	ss := strings.Split(path, "/")
	if len(ss) < 2 {
		err = ErrInvalidPath
		return
	}

	for i := 1; i <= 2; i++ {
		if len(ss) <= i {
			continue
		}
		j := strings.Index(ss[i], ".v")
		if j == -1 {
			continue
		}
		ver = ss[i][j+1:]
	}
	pkg = path
	return
}

func gopkgRepository(pkg string) (repo, base string, err error) {
	ss := strings.Split(pkg, "/")
	if len(ss) < 2 {
		err = ErrInvalidPackage
		return
	}

	var i int
	var owner, name string
	for i = 1; i <= 2; i++ {
		if len(ss) <= i {
			continue
		}

		name = ss[i]
		j := strings.Index(name, ".v")
		if j != -1 {
			name = name[:j]
			break
		}
	}
	if i == 1 {
		owner = "go-" + name
		base = filepath.Join(ss[:2]...)
	} else {
		owner = ss[1]
		base = filepath.Join(ss[:3]...)
	}

	repo = fmt.Sprintf("github.com/%s/%s", owner, name)
	return
}

func init() {
	match := prefixMatchFunc("gopkg.in")
	repository.Registerf(match, gopkgRepository)
	parser.Registerf(match, gopkgParse)
}
