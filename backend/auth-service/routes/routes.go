package routes

import (
	"net/http"

	"github.com/SergeyBarshin/atom-hack-2025/backend/auth-service/auth"
)

func SetupRoutes() {
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	// http.HandleFunc("/get-uuid", auth.GetUUIDHandler)
	http.HandleFunc("/logout", auth.LogoutUser)
}
