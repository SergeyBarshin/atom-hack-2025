package interfaces

import (
	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/models"
)

// методы для работы с таблицей graph
type GraphRepository interface {
	Create(graph *models.Graph) error
	Get(userGuid string, graphNum int) (*models.Graph, error)
	Update(graph *models.Graph) error
	Delete(userGuid string, graphNum int) error
}

// для работы с сервисом
type DataService interface {
	CreateGraph(graph *models.Graph) error
	GetGraph(userGuid string, graphNum int) (*models.Graph, error)
	UpdateGraph(retort *models.Graph) error
	DeleteGraph(userGuid string, graphNum int) error
}
