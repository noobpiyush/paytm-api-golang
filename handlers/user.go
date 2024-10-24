package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/noobpiyush/paytm-api/jwt"
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
	// Method checking

	if r.Method != http.MethodPost {
		sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	//content typpe cehecking

	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	// read body with size limit of 1 mb

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	log.Print(string(body))

	//parse json

	var signupReq SignupRequest

	if err := json.Unmarshal(body, &signupReq); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	if signupReq.FullName == "" || signupReq.Email == "" || signupReq.Password == "" {
		sendError(w, http.StatusBadRequest, "fullName, email and password are required")
		return
	}

	log.Printf("Signup request received for user: %s", signupReq.FullName)

	token, errFromToken := jwt.CreateToken(signupReq.FullName)

	if errFromToken != nil {
		log.Fatalf("Failed to create token: %v", errFromToken)
		return
	}
	log.Printf("printing ip \n")
	GetIP(r)

	if len(token) == 0 {
		log.Fatal("error creating token")
		return
	}

	sendSuccess(w, http.StatusCreated, "user registered successfully", map[string]string{
		"fullName": signupReq.FullName,
		"email":    signupReq.Email,
		"token":    token,
	})

}

func GetIP(r *http.Request) string {
	// Get IP from X-FORWARDED-FOR header
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	// Method checking

	if r.Method != http.MethodPost {
		sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	//content typpe cehecking

	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	// read body with size limit of 1 mb

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		sendError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// log.Print(string(body))

	//parse json

	var signinReq SigninRequest

	if err := json.Unmarshal(body, &signinReq); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}
	//validate fields
	if signinReq.Email == "" || signinReq.Password == "" {
		sendError(w, http.StatusBadRequest, "email or password is wrong")
		return
	}

	//bunch of db calls

	token, err := jwt.CreateToken(signinReq.Email)

	if err != nil {
		log.Printf("failed to create token %v", err)
		sendError(w, http.StatusInternalServerError, "error creating token")
		return
	}

	clientIP := GetIP(r)

	log.Printf("Signin attempt from IP: %s for email: %s\n", clientIP, signinReq.Email)

	// 9. Send success response with token
	sendSuccess(w, http.StatusOK, "signin successful", map[string]string{
		"email": signinReq.Email,
		"token": token,
	})

}
