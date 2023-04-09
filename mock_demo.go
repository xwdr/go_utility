package utils

// Service .
type Service struct {
	api Api
}

// FindUser .
func (s *Service) FindUser(id int64) (string, error) {
	return s.api.GetNameById(id)
}
