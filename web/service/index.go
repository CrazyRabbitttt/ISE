package service

import (
	"Search-Engine/search-engine/container"
	"Search-Engine/search-engine/model"
)

type IndexService struct {
	container *container.Container
}

func NewIndexService() *IndexService {
	return &IndexService{
		container: container.GlobalContainer,
	}
}

func (s *IndexService) AddIndexDoc(doc *model.IndexDoc) error {
	// 往Engine对应的chan中添加数据
	s.container.GetEngine().AddIndexDoc2Chan(doc)
	return nil
}
