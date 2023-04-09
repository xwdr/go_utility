package utils

//go:generate mockgen -destination ./demo_mock.go -package utils -source mock_interface.go
// mockgen -source 需要mock的文件名 -destination 生成的mock文件名 -package 生成mock文件的包名
// go install github.com/golang/mock/mockgen@latest
// Search .
type Api interface {
	GetNameById(id int64) (string, error)
}

// var _ Api = new(TestService)

// service.service.go
type TestService struct{}

// NewTestService .
func NewTestService() Api {
	return &TestService{}
}

// GetNameById 找到用户
func (t *TestService) GetNameById(id int64) (string, error) {
	return "this is test", nil
}
