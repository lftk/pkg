package cmd

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/4396/pkg/conf"
	"github.com/4396/pkg/dep"
	"github.com/4396/pkg/store"
	"github.com/4396/pkg/vcs"
	"github.com/urfave/cli"
)

var (
	Get = cli.Command{
		Name:  "get",
		Usage: "fetch package",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "l",
				Usage: "local verdor",
			},
			cli.BoolFlag{
				Name:  "u",
				Usage: "update dependencies",
			},
			cli.BoolFlag{
				Name:  "v",
				Usage: "verbose progress",
			},
		},
		Action: func(c *cli.Context) error {
			return fetchPackages(c.Args(), c.Bool("l"), c.Bool("u"), c.Bool("v"))
		},
	}
)

func fetchPackages(paths []string, local, update, verbose bool) (err error) {
	var dir string
	if local {
		dir, err = os.Getwd()
		if err != nil {
			return
		}
		dir = filepath.Join(dir, ".vendor")
	} else {
		path := os.Getenv("GOPATH")
		ss := filepath.SplitList(path)
		if len(ss) > 0 {
			dir = ss[0]
		}
	}
	os.MkdirAll(dir, os.ModePerm)

	dones := make(map[string]string)
	for _, path := range paths {
		err = fetchPackage(dir, path, dones, update, verbose)
		if err != nil {
			return
		}
	}

	if local {
		src := filepath.Join(dir, "src")
		dst := filepath.Join(filepath.Dir(dir), "vendor")
		for path := range dones {
			oldpath := filepath.Join(src, path)
			newpath := filepath.Join(dst, path)
			os.MkdirAll(filepath.Dir(newpath), os.ModePerm)
			os.RemoveAll(newpath)
			err = os.Rename(oldpath, newpath)
			if err != nil {
				return
			}
		}
		os.RemoveAll(dir)
	}
	return
}

type actionError struct {
	Action string
	Err    error
}

func warpError(action string, err error) error {
	return &actionError{action, err}
}

func (e actionError) Error() string {
	return fmt.Sprintf("%s %s", e.Action, e.Err.Error())
}

func fetchPackage(dir, path string, dones map[string]string, update, verbose bool) (err error) {
	logger(verbose).Log(path)
	pkg, ver, err := vcs.Parse(path)
	if err != nil {
		err = warpError("vcs.Parse", err)
		return
	}
	logger(verbose).Log("  [package]%s", pkg)
	if ver != "" {
		logger(verbose).Log("  [version]%s", ver)
	}

	repo, base, err := vcs.Repository(pkg)
	if err != nil {
		err = warpError("vcs.Repository", err)
		return
	}
	logger(verbose).Log("  [base]%s", base)
	logger(verbose).Log("  [repository]%s", repo)

	sha, ok := dones[base]
	if ok {
		logger(verbose).Log("  [revision]%s", sha)
		return
	}

	sha, err = revision(repo, ver)
	if err != nil {
		err = warpError("vcs.Revision", err)
		return
	}
	logger(verbose).Log("  [revision]%s", sha)

	file, err := archive(repo, sha)
	if err != nil {
		err = warpError("vcs.Archive", err)
		return
	}

	dst := filepath.Join(dir, "/src", base)
	switch filepath.Ext(file) {
	case ".zip":
		err = unzip(file, dst)
		if err != nil {
			err = warpError("unzip", err)
			return
		}
	default:
		err = errors.New("unsupported ext")
		return
	}

	dones[base] = sha
	if !update {
		return
	}

	imports, err := dep.Imports(pkg, base, dir, false)
	if err != nil {
		err = warpError("dep.Imports", err)
		return
	}

	for _, path := range imports {
		err = fetchPackage(dir, path, dones, update, verbose)
		if err != nil {
			return
		}
	}
	return
}

type logger bool

func (l logger) Log(format string, v ...interface{}) {
	if l {
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		name := strings.Replace(f.Name, "/", string(filepath.Separator), -1)
		ss := strings.SplitN(name, string(filepath.Separator), 2)
		path := filepath.Join(dst, ss[len(ss)-1])
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer w.Close()

		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func httpGet(path string) (resp *http.Response, err error) {
	registry, err := conf.Registry("")
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, registry+path, nil)
	if err != nil {
		return
	}

	token, err := conf.Token("")
	if err != nil {
		return
	}

	req.Header.Set("pkg-access-token", token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}
	return
}

func revision(repo, ver string) (sha string, err error) {
	resp, err := httpGet(fmt.Sprintf("/v1/revision?repo=%s&ver=%s", repo, ver))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var v struct {
		SHA string `json:"sha"`
	}
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return
	}

	sha = v.SHA
	return
}

func archive(repo, sha string) (file string, err error) {
	file, ok := store.Get(repo, sha)
	if ok {
		return
	}

	resp, err := httpGet(fmt.Sprintf("/v1/archive?repo=%s&sha=%s", repo, sha))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	filename := resp.Header.Get("Content-Disposition")
	err = store.Put(repo, sha, filepath.Ext(filename), resp.Body)
	if err != nil {
		return
	}

	file, _ = store.Get(repo, sha)
	return
}
