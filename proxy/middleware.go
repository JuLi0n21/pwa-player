package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_cookie")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := GetUserByCookie(cookie.Value)
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_cookie")
		if err != nil || cookie.Value == "" {

			newCookie := &http.Cookie{
				Name:     "session_cookie",
				Value:    "session_cookie_" + generateRandomString(64),
				Expires:  time.Now().Add(time.Hour * 10000),
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
			}

			http.SetCookie(w, newCookie)
		}

		ctx := context.WithValue(r.Context(), "cookie", cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Print(err)
		return "IF_U_SEE_THIS_UR_THE_GOAT"
	}

	return base64.URLEncoding.EncodeToString(bytes)
}
