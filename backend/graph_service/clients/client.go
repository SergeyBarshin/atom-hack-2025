package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
	url := os.Getenv("HTTP_URL")

	if graph == nil {
		return fmt.Errorf("graph is nil")
	}

	jsonData, err := json.Marshal(graph.GraphData)
	if err != nil {
		return err
	}
	// PUT запрос на обновление данных
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned: %d %s",
			resp.StatusCode,
			http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var newGraphData map[string]interface{}
	if err := json.Unmarshal(body, &newGraphData); err != nil {
		return err
	}
	graph.GraphData = newGraphData

	return nil
}
