package vcs

func init() {
	for k, v := range map[string]string{
		"google.golang.org/appengine": "github.com/golang/appengine",
		"google.golang.org/grpc":      "github.com/grpc/grpc-go",
		"google.golang.org/genproto":  "github.com/google/go-genproto",
	} {
		Proxy(k, v, k)
	}
}
