package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/interfaces"
	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/models"
)

// реализация GraphRepository для PostgreSQL
type GraphRepositoryPostgres struct {
	db *sql.DB
}

// конструктор для GraphRepositoryPostgres
func NewGraphRepositoryPostgres(db *sql.DB) interfaces.GraphRepository {
	return &GraphRepositoryPostgres{db: db}
}

func (r *GraphRepositoryPostgres) Create(graph *models.Graph) error { //работает
	// Реализация создания графа в Postgres

	log.Printf("smt2")
	query := `
        INSERT INTO graphs (user_guid, graph_num, graph_data)
        VALUES ($1, $2, $3)
    `
	graphDataJSON, err := json.Marshal(graph.GraphData)
	if err != nil {
		return fmt.Errorf("failed to marshal graph data to JSON: %w", err)
	}

	log.Printf("smt3")

	_, err = r.db.Exec(query, graph.UserGUID, graph.GraphNum, graphDataJSON)
	return err
}

func (r *GraphRepositoryPostgres) Get(userGuid string, graphNum int) (*models.Graph, error) { //работает
	//получение графа по конкретному юзеру и конкретному номеру графа
	query := `
        SELECT user_guid, graph_num, graph_data
        FROM graphs
        WHERE user_guid = $1 AND graph_num = $2
    `
	row := r.db.QueryRow(query, userGuid, graphNum)
	graph := &models.Graph{}
	var graphDataJSON []byte
	err := row.Scan(&graph.UserGUID, &graph.GraphNum, &graphDataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("graph not found")
		}
		return nil, err
	}

	err = json.Unmarshal(graphDataJSON, &graph.GraphData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph data from JSON: %w", err)
	}

	return graph, nil
}

func (r *GraphRepositoryPostgres) Update(graph *models.Graph) error { //работает
	//обновление графа
	query := `
        UPDATE graphs
        SET graph_data = $3
        WHERE user_guid = $1 AND graph_num = $2
    `
	graphDataJSON, err := json.Marshal(graph.GraphData)
	if err != nil {
		return fmt.Errorf("failed to marshal graph data to JSON: %w", err)
	}

	_, err = r.db.Exec(query, graph.UserGUID, graph.GraphNum, graphDataJSON)
	return err
}

func (r *GraphRepositoryPostgres) Delete(userGuid string, graphNum int) error { //работает
	//удаление графа
	query := `
        DELETE FROM graphs
        WHERE user_guid = $1 AND graph_num = $2
    `
	result, err := r.db.Exec(query, userGuid, graphNum)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	// Проверяем, была ли удалена хотя бы одна строка
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Возвращаем стандартную ошибку "нет данных"
	}

	return nil
}

func (r *GraphRepositoryPostgres) List(userGuid string) ([]*models.Graph, error) {
	query := `
        SELECT user_guid, graph_num, graph_data
        FROM graphs
        WHERE user_guid = $1
        ORDER BY graph_num  // Для упорядочивания по номеру графа
    `
	rows, err := r.db.Query(query, userGuid)
	if err != nil {
		return nil, fmt.Errorf("failed to query graphs: %w", err)
	}
	defer rows.Close()

	var graphs []*models.Graph
	for rows.Next() {
		var graph models.Graph
		var graphDataJSON []byte

		err := rows.Scan(&graph.UserGUID, &graph.GraphNum, &graphDataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan graph row: %w", err)
		}

		if err := json.Unmarshal(graphDataJSON, &graph.GraphData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal graph data: %w", err)
		}

		graphs = append(graphs, &graph)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return graphs, nil
}
