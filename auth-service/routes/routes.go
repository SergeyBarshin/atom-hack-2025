package routes

import (
	"net/http"

	"auth-service/auth"
)

func SetupRoutes() {
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	// http.HandleFunc("/get-uuid", auth.GetUUIDHandler)
	http.HandleFunc("/logout", auth.LogoutUser)
}
