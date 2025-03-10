package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/grpc"
	"github.com/SergeyBarshin/atom-hack-2025/backend/graph_service/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func CheckJWT(c *gin.Context) (string, error) {

	//соединение с gRPC сервером
	port := os.Getenv("GRPC_PORT")

	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	jwtToken := c.Request.Header.Get("Authorization")

	// Отправляем запрос с токеном
	req := &pb.TokenRequest{Token: jwtToken}
	resp, err := client.GetUUID(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.Uuid, nil
}

func UpdateData(graph *models.Graph) error {
	if graph == nil {
        return fmt.Errorf("graph is nil")
    }

    url := os.Getenv("HTTP_URL")
    if url == "" {
        return fmt.Errorf("HTTP_URL environment variable is not set")
    }

    log.Printf("URL: %s", url) // Логируем URL

    jsonData, err := json.Marshal(graph.GraphData)
    if err != nil {
        return fmt.Errorf("ошибка маршалинга JSON: %w", err)
    }
    log.Printf("Тело запроса: %s", string(jsonData)) // Логируем тело

    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("ошибка создания запроса: %w", err)
    }
    req.Header.Set("Content-Type", "application/json; charset=utf-8")

    // Клиент с таймаутом и логированием
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    log.Printf("Отправка запроса...")
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Фатальная ошибка при отправке: %v", err) // Детальный лог
        return fmt.Errorf("ошибка запроса: %w", err)
    }
    defer resp.Body.Close()

    log.Printf("Статус ответа: %d", resp.StatusCode) // Лог 

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned: %d %s",
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
	}

	log.Printf("5")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("6")

	var newGraphData map[string]interface{}
	if err := json.Unmarshal(body, &newGraphData); err != nil {
		return err
	}
	graph.GraphData = newGraphData

	log.Printf("7")
	return nil
}
