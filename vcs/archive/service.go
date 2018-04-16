package archive

type Service interface {
	Match(repo string) (ok bool)
	Archive(repo, sha string) (url string, err error)
}

type MatchFunc func(repo string) (ok bool)
type ArchiveFunc func(repo, sha string) (url string, err error)

type service struct {
	match   MatchFunc
	archive ArchiveFunc
}

func (s *service) Match(repo string) bool {
	return s.match(repo)
}

func (s *service) Archive(repo, sha string) (string, error) {
	return s.archive(repo, sha)
}

func Make(match MatchFunc, archive ArchiveFunc) Service {
	return &service{match, archive}
}

var (
	services []Service
)

func Register(service Service) {
	services = append(services, service)
}

func Registerf(match MatchFunc, archive ArchiveFunc) {
	services = append(services, &service{match, archive})
}

func Select(repo string) (service Service, ok bool) {
	for _, v := range services {
		ok = v.Match(repo)
		if ok {
			service = v
			return
		}
	}
	return
}
