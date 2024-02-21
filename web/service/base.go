package service

import (
	"Search-Engine/search-engine/container"
	"Search-Engine/search-engine/model"
)

type BaseService struct {
	container *container.Container
}

func NewBaseService() *BaseService {
	return &BaseService{
		container: container.GlobalContainer,
	}
}

func (s *BaseService) Query(searchRequest *model.SearchRequest) (*model.SimpleSearchResponse, error) {
	return container.GlobalContainer.GetEngine().Search(searchRequest)
}

func (s *BaseService) SearchRemind(query string) ([]string, error) {
	return container.GlobalContainer.GetEngine().SearchRemind(query)
}

func (s *BaseService) InitReminder(querys []string) {
	container.GlobalContainer.GetEngine().InitReminder(querys)
}
