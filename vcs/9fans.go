package vcs

func init() {
	for k, v := range map[string]string{
		"9fans.net/go/acme":  "github.com/9fans/go",
		"9fans.net/go/draw":  "github.com/9fans/go",
		"9fans.net/go/games": "github.com/9fans/go",
		"9fans.net/go/plan9": "github.com/9fans/go",
		"9fans.net/go/plumb": "github.com/9fans/go",
	} {
		Proxy(k, v, "9fans.net/go")
	}
}
