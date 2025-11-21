package app

// Service defines the main service's structure with all dependencies in it.
type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s *Service) Start() {
}
