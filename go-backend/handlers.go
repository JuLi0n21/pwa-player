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

	_ "backend/docs"

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

func (s *Server) registerRoutes() {
	http.HandleFunc("/api/v1/ping", s.ping)
	http.HandleFunc("/api/v1/login", s.login)

	http.HandleFunc("/api/v1/song/{hash}/", s.song)
	http.HandleFunc("/api/v1/songs/recents", s.recents)
	http.HandleFunc("/api/v1/songs/favorites", s.favorites)

	http.HandleFunc("/api/v1/collection", s.collection)
	http.HandleFunc("/api/v1/search/collections", s.collectionSearch)

	http.HandleFunc("/api/v1/search/active", s.activeSearch)
	http.HandleFunc("/api/v1/search/artist", s.artistSearch)

	http.HandleFunc("/api/v1/audio/{filepath}", s.songFile)
	http.HandleFunc("/api/v1/image/{filepath}", s.imageFile)
	http.Handle("/swagger/", httpSwagger.WrapHandler)
}

func run(s *Server) {
	s.registerRoutes()
	fmt.Println("starting server on http://localhost" + s.Port)
	log.Fatal(http.ListenAndServe(s.Port, nil))
}

// ping godoc
//
//	@Summary		Check server health
//	@Description	Returns a pong response if the server is running
//	@Tags			health
//	@Success		200	{string}	string	"pong"
//	@Router			/ping [get]
func (s *Server) ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

// login godoc
//
//	@Summary		Redirect to login page
//	@Description	Redirects users to an external authentication page
//	@Tags			auth
//	@Success		307	{string}	string	"Temporary Redirect"
//	@Router			/login [get]
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://proxy.illegalesachen.download/login", http.StatusTemporaryRedirect)
}

// song godoc
//
//	@Summary		Get a song by its hash
//	@Description	Retrieves a song using its unique hash identifier.
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			hash	path		string	true	"Song hash"
//	@Success		200		{object}	Song
//	@Failure		400		{string}	string	"Invalid parameter"
//	@Failure		404		{string}	string	"Song not found"
//	@Router			/song/{hash} [get]
func (s *Server) song(w http.ResponseWriter, r *http.Request) {

	hash := r.PathValue("hash")
	if hash == "" {
		http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
		return
	}

	song, err := getSong(s.Db, hash)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "beatmap not found by hash", http.StatusNotFound)
		return
	}

	writeJSON(w, song, http.StatusOK)
}

// recents godoc
//
//	@Summary		Get a list of recent songs
//	@Description	Retrieves recent songs with pagination support.
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"		default(10)
//	@Param			offset	query		int	false	"Offset"	default(0)
//	@Success		200		{array}		Song
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/songs/recents [get]
func (s *Server) recents(w http.ResponseWriter, r *http.Request) {

	limit, offset := pagination(r)
	recent, err := getRecent(s.Db, limit, offset)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, recent, http.StatusOK)
}

// favorites godoc
//
//	@Summary		Get a list of favorite songs based on a query
//	@Description	Retrieves favorite songs filtered by a query with pagination support.
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	true	"Search query"
//	@Param			limit	query		int		false	"Limit"		default(10)
//	@Param			offset	query		int		false	"Offset"	default(0)
//	@Success		200		{array}		Song
//	@Failure		400		{string}	string	"Invalid parameter"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/songs/favorites [get]
func (s *Server) favorites(w http.ResponseWriter, r *http.Request) {
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
}

// collection godoc
//
//	@Summary		Get a collection of songs by index
//	@Description	Retrieves a collection of songs using the provided index.
//	@Tags			songs
//	@Accept			json
//	@Produce		json
//	@Param			index	query		int	false	"Index"
//	@Param			name	query		string 	false	"Index"
//	@Success		200		{array}		Song
//	@Failure		400		{string}	string	"Invalid parameter"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/collection [get]
func (s *Server) collection(w http.ResponseWriter, r *http.Request) {
	limit, offset := pagination(r)
	name := r.URL.Query().Get("name")
	if name != "" {
		c, err := getCollectionByName(s.Db, limit, offset, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, c, http.StatusOK)
		return
	}

	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
		return
	} else {
		c, err := getCollection(s.Db, limit, offset, index)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, c, http.StatusOK)
		return
	}

}

// @Summary		Searches collections based on a query
// @Description	Searches collections in the database based on the query parameter
// @Tags			search
// @Accept			json
// @Produce		json
// @Param			query	query		string		true	"Search query"
// @Param			limit	query		int			false	"Limit the number of results"	default(10)
// @Param			offset	query		int			false	"Offset for pagination"			default(0)
// @Success		200		{array}		Collection	"List of collections"
// @Failure		400		{object}	string		"Bad Request"
// @Failure		500		{object}	string		"Internal Server Error"
// @Router			/search/collections [get]
func (s *Server) collectionSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")

	limit, offset := pagination(r)

	preview, err := getCollections(s.Db, q, limit, offset)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, preview, http.StatusOK)
}

// @Summary		Searches active records based on a query
// @Description	Searches active records in the database based on the query parameter
// @Tags			search
// @Accept			json
// @Produce		json
// @Param			query	query		string			true	"Search query"
// @Param			limit	query		int				false	"Limit the number of results"	default(10)
// @Param			offset	query		int				false	"Offset for pagination"			default(0)
// @Success		200		{object}	ActiveSearch	"Active search result"
// @Failure		400		{object}	string			"Bad Request"
// @Failure		500		{object}	string			"Internal Server Error"
// @Router			/search/active [get]
func (s *Server) activeSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	if q == "" {
		http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
		return
	}

	//TODO
	limit, offset := pagination(r)

	recent, err := getSearch(s.Db, q, limit, offset)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, recent, http.StatusOK)
}

// @Summary		Searches for artists based on a query
// @Description	Searches for artists in the database based on the query parameter
// @Tags			search
// @Accept			json
// @Produce		json
// @Param			query	query		string	true	"Search query"
// @Param			limit	query		int		false	"Limit the number of results"	default(10)
// @Param			offset	query		int		false	"Offset for pagination"			default(0)
// @Success		200		{array}		string	"List of artists"
// @Failure		400		{object}	string	"Bad Request"
// @Failure		500		{object}	string	"Internal Server Error"
// @Router			/search/artist [get]
func (s *Server) artistSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("query")
	if q == "" {
		http.Error(w, ErrRequiredParameterNotPresent.Error(), http.StatusBadRequest)
		return
	}

	//TODO
	limit, offset := pagination(r)

	a, err := getArtists(s.Db, q, limit, offset)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, a, http.StatusOK)
}

// @Summary		Retrieves a song file by its encoded path
// @Description	Retrieves a song file from the server based on the provided encoded filepath
// @Tags			files
// @Accept			json
// @Produce		json
// @Param			filepath	path		string	true	"Base64 encoded file path"
// @Success		200			{file}		File	"The requested song file"
// @Failure		400			{object}	string	"Bad Request"
// @Failure		404			{object}	string	"File Not Found"
// @Router			/audio/{filepath} [get]
func (s *Server) songFile(w http.ResponseWriter, r *http.Request) {
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
}

// @Summary		Retrieves an image file by its encoded path
// @Description	Retrieves an image file from the server based on the provided encoded filepath
// @Tags			files
// @Accept			json
// @Produce		json
// @Param			filepath	path		string	true	"Base64 encoded file path"
// @Success		200			{file}		File	"The requested image file"
// @Failure		400			{object}	string	"Bad Request"
// @Failure		404			{object}	string	"File Not Found"
// @Router			/image/{filepath} [get]
func (s *Server) imageFile(w http.ResponseWriter, r *http.Request) {
	f := r.PathValue("filepath")
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
