package revision

type Service interface {
	Match(repo string) (ok bool)
	Revision(repo, ver string) (sha string, err error)
}

type MatchFunc func(repo string) (ok bool)
type RevisionFunc func(repo, ver string) (sha string, err error)

type service struct {
	match    MatchFunc
	revision RevisionFunc
}

func (s *service) Match(repo string) bool {
	return s.match(repo)
}

func (s *service) Revision(repo, ver string) (string, error) {
	return s.revision(repo, ver)
}

func Make(match MatchFunc, revision RevisionFunc) Service {
	return &service{match, revision}
}

var (
	services []Service
)

func Register(service Service) {
	services = append(services, service)
}

func Registerf(match MatchFunc, revision RevisionFunc) {
	services = append(services, &service{match, revision})
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
