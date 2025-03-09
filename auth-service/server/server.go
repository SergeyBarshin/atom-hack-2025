package server

import (
	"context"
	"errors"

	"auth-service/auth"
	"auth-service/pb"

	"github.com/golang-jwt/jwt/v5"
)

// AuthServer реализует gRPC сервис AuthService
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
}

// GetUUID обрабатывает gRPC-запрос с токеном и возвращает UUID
func (s *AuthServer) GetUUID(ctx context.Context, req *pb.TokenRequest) (*pb.UUIDResponse, error) {
	// Проверяем и декодируем JWT
	token, err := auth.ValidateJWT(req.Token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Проверяем, является ли claims типом jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Извлекаем UUID из claims
	userUUID, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.New("uuid not found in token")
	}

	// Возвращаем UUID
	return &pb.UUIDResponse{Uuid: userUUID}, nil
}
