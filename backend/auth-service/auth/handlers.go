package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"database/sql"

	"github.com/SergeyBarshin/atom-hack-2025/backend/auth-service/db"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"pass"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Регистрация
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email string `json:"email"`
		Pass  string `json:"pass"`
	}

	log.Printf("smt1")

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	log.Printf("smt2")

	creds.Email = strings.ToLower(strings.TrimSpace(creds.Email))
	if creds.Email == "" || creds.Pass == "" {
		http.Error(w, `{"error": "Email and password are required"}`, http.StatusBadRequest)
		return
	}

	log.Printf("smt3")

	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", creds.Email).Scan(&exists)
	if err != nil {
		log.Printf("Check user exists error: %v", err)
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, `{"error": "User already exists"}`, http.StatusConflict)
		return
	}


	log.Printf("smt4")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Pass), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hash error: %v", err)
		http.Error(w, `{"error": "Error hashing password"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("smt5")

	var userUUID string
	err = db.DB.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING uuid", creds.Email, hashedPassword).Scan(&userUUID)
	if err != nil {
		log.Printf("Insert user error: %v", err)
		http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("smt6")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"uuid":    userUUID,
	})

	log.Printf("smt7")
}

// Авторизация
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Проверяем, есть ли пользователь в БД
	var userUUID, hashedPassword string
	err = db.DB.QueryRow("SELECT uuid, password FROM users WHERE email=$1", req.Email).Scan(&userUUID, &hashedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 2. Проверяем пароль (сравниваем хеш)
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 3. Генерируем JWT
	token, err := GenerateJWT(userUUID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// 4. Отправляем токен в куки
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful, token set in cookie"))
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Токен сразу удаляется
	})

	w.Write([]byte("Logged out"))
}
