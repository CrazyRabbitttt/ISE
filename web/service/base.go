package service

import "Search-Engine/search-engine/container"

type BaseService struct {
	container *container.Container
}

func NewBaseService() *BaseService {
	return &BaseService{
		container: container.GlobalContainer,
	}
}
