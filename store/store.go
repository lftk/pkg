package store

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	dir   string
	mtx   sync.RWMutex
	repos map[string]archives
)

type archives map[string]string

func init() {
	var home string
	if runtime.GOOS == "windows" {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else {
		home = os.Getenv("HOME")
	}

	_, err := Dir(filepath.Join(home, ".gopkg"))
	if err != nil {
		log.Fatal(err)
	}
}

func walk(dir string) (repos map[string]archives, err error) {
	repos = make(map[string]archives)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		base := filepath.Base(rel)
		ss := strings.Split(base, ".")
		if len(ss) != 2 {
			return nil
		}

		sha := ss[0]
		ext := "." + ss[1]
		repo := filepath.Dir(rel)
		archives, ok := repos[repo]
		if !ok {
			archives = make(map[string]string)
			repos[repo] = archives
		}
		archives[sha] = ext
		return nil
	})
	return
}

func Dir(path string) (string, error) {
	if path != "" && dir != path {
		os.MkdirAll(path, os.ModePerm)
		r, err := walk(path)
		if err != nil {
			return "", err
		}

		mtx.Lock()
		repos = r
		dir = path
		mtx.Unlock()
	}
	return dir, nil
}

func Map(f func(repo, sha, ext string) error) (err error) {
	mtx.RLock()
	defer mtx.RUnlock()

	for repo, archives := range repos {
		for sha, ext := range archives {
			err = f(repo, sha, ext)
			if err != nil {
				return
			}
		}
	}
	return
}

func Get(repo, sha string) (path string, ok bool) {
	mtx.RLock()
	defer mtx.RUnlock()

	archives, ok := repos[repo]
	if !ok {
		return
	}

	ext, ok := archives[sha]
	if !ok {
		return
	}

	path = filepath.Join(dir, repo, sha) + ext
	return
}

func Put(repo, sha, ext string, r io.Reader) (err error) {
	path := filepath.Join(dir, repo)
	os.MkdirAll(path, os.ModePerm)

	prefix := sha + ext + "."
	f, err := ioutil.TempFile(path, prefix)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return
	}

	path = filepath.Join(dir, repo, sha+ext)
	err = os.Rename(f.Name(), path)
	if err != nil {
		return
	}

	mtx.Lock()
	archives, ok := repos[repo]
	if !ok {
		archives = make(map[string]string)
		repos[repo] = archives
	}
	archives[sha] = ext
	mtx.Unlock()
	return
}
