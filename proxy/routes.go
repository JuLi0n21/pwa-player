package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func run() error {
	mux := http.NewServeMux()

	mux.Handle("/me", AuthMiddleware(http.HandlerFunc(MeHandler)))
	mux.Handle("/login", http.HandlerFunc(LoginRedirect))

	fmt.Println("Starting Server on :42000")

	//global middleware
	handler := CookieMiddleware(mux)

	return http.ListenAndServe(":42000", handler)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(*User)

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
