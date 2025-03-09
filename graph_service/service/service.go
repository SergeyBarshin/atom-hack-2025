package service

import (
	"graph_service/interfaces"
	"graph_service/models"
)

// DataService - реализация DataService
type DataService struct {
	graphRepo interfaces.GraphRepository
}

// NewDataService - конструктор для DataService
func NewDataService(graphRepo interfaces.GraphRepository) interfaces.DataService {
	return &DataService{
		graphRepo: graphRepo,
	}
}

// CreateGraph - создание графа
func (data *DataService) CreateGraph(graph *models.Graph) error {
	return data.graphRepo.Create(graph)
}

// GetGraph - получение графа
func (s *DataService) GetGraph(userGuid string, graphNum int) (*models.Graph, error) {
	return s.graphRepo.Get(userGuid, graphNum)
}

// UpdateGraph - обновление графа
func (s *DataService) UpdateGraph(graph *models.Graph) error {
	return s.graphRepo.Update(graph)
}

// DeleteGraph - удаление графа
func (s *DataService) DeleteGraph(userGuid string, graphNum int) error {
	return s.graphRepo.Delete(userGuid, graphNum)
}
