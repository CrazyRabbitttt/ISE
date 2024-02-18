package service

import "Search-Engine/search-engine/container"

type DataBaseService struct {
	container *container.Container
}

func NewDataBaseService() *DataBaseService {
	return &DataBaseService{
		container: container.GlobalContainer,
	}
}
