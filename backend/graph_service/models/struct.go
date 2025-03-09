package models

import (
	"fmt"
)

type Graph struct {
	UserGUID  string                 `db:"user_guid" json:"user_guid" binding:"required,uuid"`
	GraphNum  int                    `db:"graph_num" json:"project_id" binding:"required,min=1"`
	GraphData map[string]interface{} `db:"graph_data" json:"project_data" binding:"required"`
}

// для проверки
func (g *Graph) String() string {
	return fmt.Sprintf(
		"UserGUID: %s, GraphNum: %d, GraphData: %v",
		g.UserGUID,
		g.GraphNum,
		g.GraphData,
	)
}
