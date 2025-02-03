package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/juli0n21/go-osu-parser/parser"
)

var ErrRequiredParameterNotPresent = errors.New("required parameter missing")
var ErrFailedToParseEncoded = errors.New("invalid encoded value")
var ErrFileNotFound = errors.New("file not found")

type Server struct {
	Port   string
	OsuDir string
	Db     *sql.DB
	OsuDb  *parser.OsuDB
}

func run(s *Server) {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://proxy.illegalesachen.download/login", http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/api/v1/songs/{hash}/", func(w http.ResponseWriter, r *http.Request) {

		hash := r.PathValue("hash")
		if hash == "" {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		song, err := getSong(s.Db, hash)
		if err != nil {
			fmt.Println(err)
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusNotFound)
			return
		}

		writeJSON(w, song, http.StatusOK)
	})

	http.HandleFunc("/api/v1/songs/recent", func(w http.ResponseWriter, r *http.Request) {

		limit, offset := pagination(r)
		recent, err := getRecent(s.Db, limit, offset)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, recent, http.StatusOK)
	})

	http.HandleFunc("/api/v1/songs/favorite", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		limit, offset := pagination(r)

		favorites, err := getFavorites(s.Db, query, limit, offset)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, favorites, http.StatusOK)
	})

	http.HandleFunc("/api/v1/collections", func(w http.ResponseWriter, r *http.Request) {
		index, err := strconv.Atoi(r.URL.Query().Get("index"))
		if err != nil {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		recent, err := getCollection(s.Db, index)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, recent, http.StatusOK)
	})

	http.HandleFunc("/api/v1/collections/", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		limit, offset := pagination(r)

		recent, err := getCollections(s.Db, query, limit, offset)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, recent, http.StatusOK)
	})

	http.HandleFunc("/api/v1/audio/{filepath}", func(w http.ResponseWriter, r *http.Request) {
		f := r.PathValue("filepath")
		if f == "" {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		filename, err := base64.RawStdEncoding.DecodeString(f)
		if err != nil {
			fmt.Println(err)
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		file, err := os.Open(s.OsuDir + "Songs/" + string(filename))
		if err != nil {
			fmt.Println(err)
			http.Error(w, ErrFileNotFound.Error(), http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			http.Error(w, ErrFileNotFound.Error(), http.StatusNotFound)
			return
		}

		http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
	})

	http.HandleFunc("/api/v1/search/active", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("query")
		if q == "" {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		limit, offset := pagination(r)

		recent, err := getSearch(s.Db, q, limit, offset)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, recent, http.StatusOK)
	})

	http.HandleFunc("/api/v1/search/artist", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("query")
		if q == "" {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		limit, offset := pagination(r)

		recent, err := getArtists(s.Db, q, limit, offset)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, recent, http.StatusOK)
	})

	http.HandleFunc("/api/v1/search/collections", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("query")
		if q == "" {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		limit, offset := pagination(r)

		recent, err := getCollections(s.Db, q, limit, offset)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, recent, http.StatusOK)
	})

	http.HandleFunc("/api/v1/images/{filename}", func(w http.ResponseWriter, r *http.Request) {
		f := r.PathValue("filename")
		if f == "" {
			http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
			return
		}

		filename, err := base64.RawStdEncoding.DecodeString(f)
		if err != nil {
			fmt.Println(err)
			http.Error(w, ErrFailedToParseEncoded.Error(), http.StatusBadRequest)
			return
		}

		http.ServeFile(w, r, s.OsuDir+"Songs/"+string(filename))
	})

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("starting server on http://localhost" + s.Port)
	log.Fatal(http.ListenAndServe(s.Port, nil))
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

func pagination(r *http.Request) (int, int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 100
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}
