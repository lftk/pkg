package repository

type Service interface {
	Match(pkg string) (ok bool)
	Repository(pkg string) (repo, base string, err error)
}

type MatchFunc func(pkg string) (ok bool)
type RepositoryFunc func(pkg string) (repo, base string, err error)

type service struct {
	match      MatchFunc
	repository RepositoryFunc
}

func (s *service) Match(repo string) bool {
	return s.match(repo)
}

func (s *service) Repository(pkg string) (string, string, error) {
	return s.repository(pkg)
}

func Make(match MatchFunc, repository RepositoryFunc) Service {
	return &service{match, repository}
}

var (
	services []Service
)

func Register(service Service, front bool) {
	if front {
		services = append([]Service{service}, services...)
	} else {
		services = append(services, service)
	}
}

func Registerf(match MatchFunc, repository RepositoryFunc, front bool) {
	Register(Make(match, repository), front)
}

func Select(pkg string) (service Service, ok bool) {
	for _, v := range services {
		ok = v.Match(pkg)
		if ok {
			service = v
			return
		}
	}
	return
}
