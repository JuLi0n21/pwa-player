package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/juli0n21/go-osu-parser/parser"
)

func run(db *sql.DB, osuDb *parser.OsuDB) {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Login endpoint")
	})

	http.HandleFunc("/api/v1/songs/", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			fmt.Fprintln(w, err)
		}

		recent, err := getRecent(db, limit, offset)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, recent)
	})

	http.HandleFunc("/api/v1/songs/recent", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			fmt.Fprintln(w, err)
		}

		recent, err := getRecent(db, limit, offset)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, recent)
		fmt.Fprintln(w, "Recent songs endpoint")
	})

	http.HandleFunc("/api/v1/songs/favorite", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			fmt.Fprintln(w, err)
		}

		recent, err := getRecent(db, limit, offset)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, recent)
		fmt.Fprintln(w, "Favorite songs endpoint")
	})

	http.HandleFunc("/api/v1/collections/{index}", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			fmt.Fprintln(w, err)
		}

		recent, err := getRecent(db, limit, offset)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, recent)
		fmt.Fprintln(w, "Collections endpoint")
	})

	http.HandleFunc("/api/v1/collections/", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			fmt.Fprintln(w, err)
		}

		recent, err := getRecent(db, limit, offset)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, recent)
		fmt.Fprintln(w, "Collection with index endpoint")
	})

	http.HandleFunc("/api/v1/audio/{hash}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Audio endpoint")
	})

	http.HandleFunc("/api/v1/search/active", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		q := r.URL.Query().Get("query")
		if err != nil || q == "" {
			fmt.Fprintln(w, errors.New("'query' is required for search"))
		}

		recent, err := getSearch(db, q, limit, offset)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, recent)
		fmt.Fprintln(w, "Active search endpoint")
	})

	http.HandleFunc("/api/v1/search/artist", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		q := r.URL.Query().Get("query")
		if err != nil {
			fmt.Fprintln(w, err)
		}

		recent, err := getArtists(db, q, limit, offset)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, recent)
		fmt.Fprintln(w, "Artist search endpoint")
	})

	http.HandleFunc("/api/v1/search/collections", func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			fmt.Fprintln(w, err)
		}
		q := r.URL.Query().Get("query")
		if err != nil {
			fmt.Fprintln(w, err)
		}

		recent, err := getCollections(db, q, limit, offset)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		fmt.Fprintln(w, recent)
	})

	http.HandleFunc("/api/v1/images/{filename}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Images endpoint")
	})

	fmt.Println("starting server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
