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

	mux.Handle("/me", AuthMiddleware(http.HandlerFunc(MeHandler)))
	mux.Handle("/login", http.HandlerFunc(LoginRedirect))

	mux.Handle("/oauth/code", http.HandlerFunc(Oauth))

	fmt.Println("Starting Server on", port)

	//global middleware
	handler := CookieMiddleware(mux)

	return http.ListenAndServe(port, handler)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(*User)

	//mask token...
	user.Token = Token{}

	w.Header().Set("Content-Type", "application/Json")

	JSONResponse(w, http.StatusOK, user)
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
	}
}
