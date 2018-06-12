package vcs

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/4396/pkg/vcs/archive"
	"github.com/4396/pkg/vcs/parser"
	"github.com/4396/pkg/vcs/repository"
	"github.com/4396/pkg/vcs/revision"
)

func githubArchive(repo, sha string) (url string, err error) {
	url = fmt.Sprintf("https://%s/archive/%s.zip", repo, sha)
	return
}

func githubParse(path string) (pkg, ver string, err error) {
	pkg = path
	i := strings.Index(path, "@")
	if i == -1 {
		return
	}

	pkg, ver = path[0:i], path[i+1:]
	return
}

func githubRepository(pkg string) (repo, base string, err error) {
	ss := strings.Split(pkg, "/")
	if len(ss) < 3 {
		err = ErrInvalidPackage
		return
	}

	repo = strings.Join(ss[:3], "/")
	base = filepath.Join(ss[:3]...)
	return
}

func githubRevision(repo, ver string) (sha string, err error) {
	url := fmt.Sprintf("https://%s.git/info/refs?service=git-upload-pack", repo)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var (
		line   []byte
		tag    = "refs/tags/" + ver
		branch = "refs/heads/" + ver
		r      = bufio.NewReader(resp.Body)
	)
	defer func() {
		if err == nil && sha == "" {
			err = ErrInvalidSHA
		}
	}()

	if ver == "" {
		for i := 0; i < 2; i++ {
			line, _, err = r.ReadLine()
			if err != nil {
				return
			}
		}

		if len(line) > 48 {
			sha = string(line[8:48])
		}
		return
	}

	for i := 0; true; i++ {
		line, _, err = r.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if i < 2 {
			continue
		}

		if !strings.Contains(string(line), tag) &&
			!strings.Contains(string(line), branch) {
			continue
		}

		if len(line) > 44 {
			sha = string(line[4:44])
		}
		return
	}
	return
}

func init() {
	match := prefixMatchFunc("github.com")
	archive.Registerf(match, githubArchive, false)
	parser.Registerf(match, githubParse, false)
	repository.Registerf(match, githubRepository, false)
	revision.Registerf(match, githubRevision, false)
}
