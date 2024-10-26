// handlers/auth.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/noobpiyush/paytm-api/db"
	"github.com/noobpiyush/paytm-api/jwt"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "only POST method is allowed")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	var signupReq SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&signupReq); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	// Validate required fields
	if signupReq.FullName == "" || signupReq.Email == "" || signupReq.Password == "" {
		sendError(w, http.StatusBadRequest, "fullName, email and password are required")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Create user
	err = db.CreateUser(signupReq.Email, string(hashedPassword), signupReq.FullName)
	if err == db.ErrUserExists {
		sendError(w, http.StatusConflict, "email already registered")
		return
	}
	if err != nil {
		log.Printf("Error creating user: %v", err)
		sendError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Generate token
	token, err := jwt.CreateToken(signupReq.Email)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		sendError(w, http.StatusInternalServerError, "error creating authentication token")
		return
	}

	sendSuccess(w, http.StatusCreated, "user registered successfully", map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"email":    signupReq.Email,
			"fullName": signupReq.FullName,
		},
	})
}

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "only POST method is allowed")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	var signinReq SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&signinReq); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	if signinReq.Email == "" || signinReq.Password == "" {
		sendError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	// Get user from database
	user, err := db.GetUserByEmail(signinReq.Email)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(signinReq.Password)); err != nil {
		sendError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Generate token
	token, err := jwt.CreateToken(user.Email)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		sendError(w, http.StatusInternalServerError, "error creating authentication token")
		return
	}

	sendSuccess(w, http.StatusOK, "signin successful", map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"email":    user.Email,
			"fullName": user.FullName,
		},
	})
}
