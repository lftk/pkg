package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"sync"

	"github.com/4396/pkg/store"
	"github.com/4396/pkg/vcs"
	"github.com/julienschmidt/httprouter"
)

var (
	addr  = flag.String("addr", ":7543", "server address")
	token = flag.String("token", "0443dbd565c01d39cb97a4e452d580986251d6c5", "access token")
)

func main() {
	flag.Parse()

	r := httprouter.New()
	r.GET("/v1/revision", auth(Revision))
	r.GET("/v1/archive", auth(Archive))

	log.Fatal(http.ListenAndServe(*addr, r))
}

func auth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if *token != r.Header.Get("pkg-access-token") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(w, r, ps)
	}
}

func Revision(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	repo := r.FormValue("repo")
	ver := r.FormValue("ver")

	sha, err := vcs.Revision(repo, ver)
	if err != nil {
		apiError(w, err)
		return
	}

	apiResult(w, map[string]string{"sha": sha})
}

func apiError(w http.ResponseWriter, err error) {
	writeJSON(w, map[string]string{"err": err.Error()})
}

func apiResult(w http.ResponseWriter, v interface{}) {
	writeJSON(w, v)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

var (
	downloading sync.Map
)

func Archive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	repo := r.FormValue("repo")
	sha := r.FormValue("sha")

	file, ok := store.Get(repo, sha)
	if ok {
		serveFile(w, r, file)
		return
	}

	errc := make(chan error)
	donec := make(chan interface{})
	key := fmt.Sprintf("%s.%s", repo, sha)
	val, loaded := downloading.LoadOrStore(key, donec)
	if loaded {
		donec = val.(chan interface{})
	} else {
		go func() {
			err := fetchArchive(repo, sha)
			if err != nil {
				errc <- err
			}
			downloading.Delete(key)
			close(donec)
		}()
	}

	select {
	case <-r.Cancel:
	case <-errc:
		w.WriteHeader(http.StatusInternalServerError)
	case <-donec:
		file, ok = store.Get(repo, sha)
		if ok {
			serveFile(w, r, file)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, file string) {
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(file))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Pragma", "public")

	http.ServeFile(w, r, file)
}

func fetchArchive(repo, sha string) (err error) {
	_, ok := store.Get(repo, sha)
	if ok {
		return
	}

	url, err := vcs.Archive(repo, sha)
	if err != nil {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = store.Put(repo, sha, path.Ext(url), resp.Body)
	return
}
