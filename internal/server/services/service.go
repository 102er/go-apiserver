package services

type BaseService struct {
	*HelloWorldService
}

var BaseServiceInstance *BaseService

func RegisterBaseService() {
	BaseServiceInstance = &BaseService{
		HelloWorldService: newHelloWorldService(),
	}
}
