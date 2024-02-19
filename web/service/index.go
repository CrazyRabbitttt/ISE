package service

import (
	"Search-Engine/search-engine/container"
	"Search-Engine/search-engine/model"
	"fmt"
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
	fmt.Println("完成IndexDoc的添加")
	return nil
}
