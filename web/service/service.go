package service

var GlobalService *Service

type Service struct {
	BaseService     *BaseService
	IndexService    *IndexService
	DataBaseService *DataBaseService
}

func InitService() {
	GlobalService = &Service{
		BaseService:     NewBaseService(),
		IndexService:    NewIndexService(),
		DataBaseService: NewDataBaseService(),
	}
}
