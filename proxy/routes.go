package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func run() error {

	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	mux.Handle("GET /me", AuthMiddleware(http.HandlerFunc(MeHandler)))
	mux.Handle("GET /login", http.HandlerFunc(LoginMiddlePage))
	mux.Handle("GET /oauth/code", http.HandlerFunc(Oauth))
	mux.Handle("POST /settings", AuthMiddleware(http.HandlerFunc(Settings)))

	//	mux.Handle("POST /setting", );

	fmt.Println("Starting Server on", port)

	//global middleware
	handler := CORS(CookieMiddleware(Logger(mux)))

	finalMux := http.NewServeMux()

	finalMux.Handle("/", handler)

	finalMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	return http.ListenAndServe(port, finalMux)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(*User)

	//mask token...
	user.Token = Token{}

	w.Header().Set("Content-Type", "application/Json")

	JSONResponse(w, http.StatusOK, user)
}

func Settings(w http.ResponseWriter, r *http.Request) {
	type settings struct {
		Sharing  *bool  `json:"sharing"`
		Endpoint string `json:"endpoint"`
	}

	user, ok := r.Context().Value("user").(*User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var s settings

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if s.Endpoint != "" {
		UpdateUserEndPoint(user.UserID, s.Endpoint)
	}

	if s.Sharing != nil {
		UpdateUserShare(user.UserID, *s.Sharing)
	}

	w.WriteHeader(http.StatusOK)
	return
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
	}
}
