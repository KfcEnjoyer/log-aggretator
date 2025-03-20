package services

import (
	"auth-service/internal/user"
	"auth-service/pkg/logger"
	"auth-service/pkg/validator"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var v = validator.NewValidator()

func Create(w http.ResponseWriter, r *http.Request, auth *AuthService) {
	defer r.Body.Close()

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ipAddress := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ipAddress = forwardedFor
	}

	decoder := json.NewDecoder(r.Body)
	var u user.RegularUser
	err := decoder.Decode(&u)
	if err != nil {
		var l logger.Log
		l.DetailedLog(
			time.Now(),
			logger.ERROR,
			"Invalid registration request format",
			logger.SECURITY,
			"system",
			"system",
			map[string]interface{}{
				"error":      err.Error(),
				"ip_address": ipAddress,
				"user_agent": r.UserAgent(),
			},
		)
		auth.LogRepo.CreateLog(l)

		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	response, err := auth.Register(u, ipAddress)
	if err != nil {
		log.Println(err)
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func LogIn(w http.ResponseWriter, r *http.Request, auth *AuthService) {
	defer r.Body.Close()

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ipAddress := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ipAddress = forwardedFor
	}

	userAgent := r.UserAgent()

	decoder := json.NewDecoder(r.Body)
	var u user.RegularUser
	err := decoder.Decode(&u)
	if err != nil {
		var l logger.Log
		l.DetailedLog(
			time.Now(),
			logger.ERROR,
			"Invalid login request format",
			logger.SECURITY,
			"system",
			"system",
			map[string]interface{}{
				"error":      err.Error(),
				"ip_address": ipAddress,
				"user_agent": userAgent,
			},
		)
		auth.LogRepo.CreateLog(l)

		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	response, ok, err := auth.LogIn(u.Username, u.Password.Plain, ipAddress, userAgent)

	if err != nil || !ok {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		http.Redirect(w, r, "/open", http.StatusFound)
		return
	}

	startTime := time.Now()

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}

	duration := time.Since(startTime).Milliseconds()

	var l logger.Log
	l.DetailedLog(
		time.Now(),
		logger.DEBUG,
		"Request processing performance",
		logger.PERFORMANCE,
		u.Username,
		response["role"],
		map[string]interface{}{
			"operation":     "login",
			"duration_ms":   duration,
			"ip_address":    ipAddress,
			"user_agent":    userAgent,
			"response_size": len(response),
		},
	)
	auth.LogRepo.CreateLog(l)
}

//func UpdateRole(w http.ResponseWriter, r *http.Request, auth *AuthService) {
//	defer r.Body.Close()
//
//	if r.Method != "POST" {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	ipAddress := r.RemoteAddr
//	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
//		ipAddress = forwardedFor
//	}
//
//	type UpdateRequest struct {
//		Username      string `json:"username"`
//		NewRole       string `json:"new_role"`
//		AdminUsername string `json:"admin_username"`
//	}
//
//	var req UpdateRequest
//	err := json.NewDecoder(r.Body).Decode(&req)
//	if err != nil {
//		http.Error(w, "Invalid request format", http.StatusBadRequest)
//		return
//	}
//
//	response, err := auth.UpdateUserRole(req.Username, req.NewRole, req.AdminUsername, ipAddress)
//	if err != nil {
//		http.Error(w, "Role update failed", http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(response)
//}

func LogSystemStartup(auth *AuthService, version, environment string) {
	var l logger.Log
	l.SystemStartup(time.Now(), version, environment)
	auth.LogRepo.CreateLog(l)
}
