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
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func CookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_cookie")
		if err != nil || cookie.Value == "" || cookie == nil {

			newCookie := &http.Cookie{
				Name:     "session_cookie",
				Value:    "session_cookie_" + generateRandomString(64),
				Expires:  time.Now().Add(time.Hour * 10000),
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
				Domain:   ".illegalesachen.download",
				SameSite: http.SameSiteLaxMode,
			}

			http.SetCookie(w, newCookie)
			cookie = newCookie
			r.AddCookie(cookie)
		}

		ctx := context.WithValue(r.Context(), "cookie", cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "https://music.illegalesachen.download")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.URL)
		next.ServeHTTP(w, r)
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
