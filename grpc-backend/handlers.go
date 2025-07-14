package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	v1 "backend/gen"
	"backend/internal/db"

	"github.com/joho/godotenv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	Env    map[string]string
	Sqlc   *db.Queries

	v1.UnimplementedMusicBackendServer
}

func (s *Server) registerRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/audio/{filepath}", s.songFile)
	mux.HandleFunc("/api/v1/image/{filepath}", s.imageFile)

	mux.HandleFunc("/api/v1/callback", s.callback)

	return mux
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingResponse, error) {
	return &v1.PingResponse{Pong: "pong"}, nil
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://proxy.illegalesachen.download/login", http.StatusTemporaryRedirect)
}

func (s *Server) Song(ctx context.Context, req *v1.SongRequest) (*v1.SongResponse, error) {

	hash := req.Hash
	if hash == "" {
		return nil, status.Errorf(codes.InvalidArgument, "hash is required and cant be empty")
	}

	song, err := getSong(s.Db, hash)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.NotFound, "beatmap not found by hash")
	}

	return &v1.SongResponse{
		Song: song.toProto(),
	}, nil
}

func (s *Server) Recent(ctx context.Context, req *v1.RecentRequest) (*v1.RecentResponse, error) {

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	recent, err := getRecent(context.Background(), s.Sqlc, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get recents")
	}

	return &v1.RecentResponse{
		Songs: toProtoSongs(recent),
	}, nil
}

func (s *Server) Favorite(ctx context.Context, req *v1.FavoriteRequest) (*v1.FavoriteResponse, error) {

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)
	favorites, err := getFavorites(s.Db, req.Query, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get favorites")
	}

	return &v1.FavoriteResponse{
		Songs: toProtoSongs(favorites),
	}, nil
}

func (s *Server) Collections(ctx context.Context, req *v1.CollectionRequest) (*v1.CollectionResponse, error) {

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	name := req.Name
	if name != "" {
		c, err := getCollectionByName(s.Db, limit, offset, name)
		if err != nil {
			fmt.Println(err)
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to fetch collection with name: %s", name))
		}
		return &v1.CollectionResponse{
			Songs: toProtoSongs(c.Songs),
			Items: int32(c.Items),
			Name:  c.Name,
		}, nil
	}
	fmt.Println(limit, offset, req.Index)
	c, err := getCollection(s.Db, limit, offset, int(req.Index))
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to fetch collection with index: %d", req.Index))
	}
	fmt.Println(c)
	return &v1.CollectionResponse{
		Songs: toProtoSongs(c.Songs),
		Items: int32(c.Items),
		Name:  c.Name,
	}, nil

}

func (s *Server) Search(ctx context.Context, req *v1.SearchSharedRequest) (*v1.SearchSharedResponse, error) {
	q := req.Query
	if q == "" {
		return nil, status.Error(codes.InvalidArgument, "query cant be empty")
	}

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	search, err := getSearch(s.Db, q, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "failed to fetch search")
	}

	return &v1.SearchSharedResponse{
		Artist: search.Artist,
		Songs:  toProtoSongs(search.Songs),
	}, nil
}

func (s *Server) SearchCollections(ctx context.Context, req *v1.SearchCollectionRequest) (*v1.SearchCollectionResponse, error) {
	q := req.Query

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	fmt.Println(req)
	preview, err := getCollections(s.Db, q, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to search for collections")
	}

	return &v1.SearchCollectionResponse{
		Collections: toProtoCollectionPreview(preview),
	}, nil
}

func (s *Server) SearchArtists(ctx context.Context, req *v1.SearchArtistRequest) (*v1.SearchArtistResponse, error) {
	q := req.Query
	if q == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	a, err := getArtists(s.Db, q, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "failed to search artists")
	}
	return &v1.SearchArtistResponse{
		Artists: toProtoArtist(a),
	}, nil
}

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

func (s *Server) callback(w http.ResponseWriter, r *http.Request) {
	cookie := r.URL.Query().Get("COOKIE")

	if cookie != "" {
		s.Env["COOKIE"] = cookie
		godotenv.Write(s.Env, ".env")
	}
}

func defaultLimit(limit int) int {
	if limit <= 0 || limit > 100 {
		return 100
	}
	return limit
}
