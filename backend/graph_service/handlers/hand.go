package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/clients"
	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/interfaces"
	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/models"

	"github.com/gin-gonic/gin"
)

// обработчик HTTP-запросов для graphs
type GraphHandler struct {
	dataService interfaces.DataService
}

// конструктор для GraphHandler
func NewGraphHandler(dataService interfaces.DataService) *GraphHandler {
	return &GraphHandler{dataService: dataService}
}

func ValidateGraphData(data map[string]interface{}) error {
	requiredKeys := []string{"matrix", "volumes", "liquid_levels", "label", "initial_liquid_level", "pressures"}
	for _, key := range requiredKeys {
		if _, ok := data[key]; !ok {
			return fmt.Errorf("the key '%s' is not in project_data", key)
		}
	}
	return nil
}

// создание графа
func (h *GraphHandler) CreateGraph(c *gin.Context) {

	/*userGUID, err := clients.CheckJWT(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}*/

	userGUID := "360e5555-e30b-20d4-a716-446655440000"

	var graph models.Graph
	graph.UserGUID = userGUID

	if err := c.ShouldBindJSON(&graph); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("smt")

	if err := clients.UpdateData(&graph); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateGraphData(graph.GraphData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.dataService.CreateGraph(&graph); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("smt4")

	c.JSON(http.StatusCreated, graph)
}

// получение графа
func (h *GraphHandler) GetGraph(c *gin.Context) {

	/*userGUID, err := clients.CheckJWT(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}*/

	userGUID := "360e5555-e30b-20d4-a716-446655440000"
	graphNum, _ := strconv.Atoi(c.Query("project_id"))

	graph, err := h.dataService.GetGraph(userGUID, graphNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, graph)
}

// обновление графа
func (h *GraphHandler) UpdateGraph(c *gin.Context) {

	/*userGUID, err := clients.CheckJWT(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}*/

	userGUID := "360e5555-e30b-20d4-a716-446655440000"

	var graph models.Graph
	graph.UserGUID = userGUID

	if err := c.ShouldBindJSON(&graph); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.dataService.UpdateGraph(&graph); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ValidateGraphData(graph.GraphData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, graph)
}

// удаление графа
func (h *GraphHandler) DeleteGraph(c *gin.Context) {

	/*userGUID, err := clients.CheckJWT(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}*/

	userGUID := "360e5555-e30b-20d4-a716-446655440000"
	graphNum, _ := strconv.Atoi(c.Query("project_id"))

	if err := h.dataService.DeleteGraph(userGUID, graphNum); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
