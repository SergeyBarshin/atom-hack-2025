package auth

import (
	"encoding/json"
	"net/http"

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

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Pass), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Вставляем в БД
	var userUUID string
	err = db.DB.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING uuid", creds.Email, hashedPassword).
		Scan(&userUUID)
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	w.Write([]byte("User registered successfully"))
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

// Проверка токена и возврат UUID
// func GetUUIDHandler(w http.ResponseWriter, r *http.Request) {
// 	authHeader := r.Header.Get("Authorization")
// 	if authHeader == "" {
// 		http.Error(w, "Missing token", http.StatusUnauthorized)
// 		return
// 	}

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 	token, err := ValidateJWT(tokenString)
// 	if err != nil || !token.Valid {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	claims, _ := token.Claims.(jwt.MapClaims)
// 	userUUID, _ := claims["uuid"].(string)

// 	w.Write([]byte(userUUID))
// }

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
