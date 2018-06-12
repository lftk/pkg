package parser

type Service interface {
	Match(path string) (ok bool)
	Parse(path string) (pkg, ver string, err error)
}

type MatchFunc func(path string) (ok bool)
type ParseFunc func(path string) (pkg, ver string, err error)

type service struct {
	match MatchFunc
	parse ParseFunc
}

func (s *service) Match(path string) bool {
	return s.match(path)
}

func (s *service) Parse(path string) (string, string, error) {
	return s.parse(path)
}

func Make(match MatchFunc, parse ParseFunc) Service {
	return &service{match, parse}
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

func Registerf(match MatchFunc, parse ParseFunc, front bool) {
	Register(Make(match, parse), front)
}

func Select(path string) (service Service, ok bool) {
	for _, v := range services {
		ok = v.Match(path)
		if ok {
			service = v
			return
		}
	}
	return
}
