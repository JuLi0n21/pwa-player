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
	"path/filepath"

	v1 "backend/gen"
	"backend/internal/db"
	sqlcdb "backend/internal/db"

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

func (s *Server) Song(ctx context.Context, req *v1.SongRequest) (*v1.SongResponse, error) {

	hash := req.Hash
	if hash == "" {
		return nil, status.Errorf(codes.InvalidArgument, "hash is required and cant be empty")
	}

	song, err := s.Sqlc.GetBeatmapByHash(ctx, sql.NullString{Valid: true, String: hash})
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.NotFound, "beatmap not found by hash")
	}

	return &v1.SongResponse{
		Song: toProtoSongSqlC(song),
	}, nil
}

func (s *Server) Recent(ctx context.Context, req *v1.RecentRequest) (*v1.RecentResponse, error) {
	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	rows, err := s.Sqlc.GetRecentBeatmaps(ctx, sqlcdb.GetRecentBeatmapsParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get recents")
	}

	return &v1.RecentResponse{
		Songs: toProtoSongsSqlC(rows),
	}, nil
}

func (s *Server) Favorite(ctx context.Context, req *v1.FavoriteRequest) (*v1.FavoriteResponse, error) {
	return nil, fmt.Errorf("not implemented!")

	/*limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)
	favorites, err := getFavorites(s.Db, req.Query, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get favorites")
	}

	return &v1.FavoriteResponse{
		Songs: toProtoSongs(favorites),
	}, nil
	*/
}

func (s *Server) Collections(ctx context.Context, req *v1.CollectionRequest) (*v1.CollectionResponse, error) {

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	name := req.Name
	if name != "" {
		c, err := s.Sqlc.GetCollectionByName(ctx,
			sqlcdb.GetCollectionByNameParams{
				Name:   sql.NullString{Valid: true, String: name},
				Limit:  int64(limit),
				Offset: int64(offset),
			},
		)
		if err != nil {
			fmt.Println(err)
			return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to fetch collection with name: %s", name))
		}

		return toProtoCollectionSqlc(c), nil
	}

	c, err := s.Sqlc.GetCollectionByOffset(ctx, sqlcdb.GetCollectionByOffsetParams{
		Offset: int64(offset),
		Limit:  int64(limit),
		Index:  int64(req.Index),
	})
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to fetch collection with index: %d", req.Index))
	}

	return toProtoCollectionoffsetSqlc(c), nil

}

func (s *Server) Search(ctx context.Context, req *v1.SearchSharedRequest) (*v1.SearchSharedResponse, error) {
	q := req.Query
	if q == "" {
		return nil, status.Error(codes.InvalidArgument, "query cant be empty")
	}
	q = "%" + q + "%"

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	search, err := s.Sqlc.SearchBeatmaps(ctx, sqlcdb.SearchBeatmapsParams{
		Title:  sql.NullString{String: q, Valid: true},
		Artist: sql.NullString{String: q, Valid: true},
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "failed to fetch search")
	}

	return &v1.SearchSharedResponse{
		Artist: "",
		Songs:  toProtoSongsSqlC(search),
	}, nil
}

func (s *Server) SearchCollections(ctx context.Context, req *v1.SearchCollectionRequest) (*v1.SearchCollectionResponse, error) {
	q := req.Query
	q = "%" + q + "%"
	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	preview, err := s.Sqlc.SearchCollection(ctx, sqlcdb.SearchCollectionParams{
		Name:   sql.NullString{String: q, Valid: true},
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to search for collections")
	}
	return toProtoCollectionsSearchPreview(preview), nil
}

func (s *Server) SearchArtists(ctx context.Context, req *v1.SearchArtistRequest) (*v1.SearchArtistResponse, error) {
	q := req.Query
	if q == "" {
		return nil, status.Error(codes.InvalidArgument, "query is required")
	}

	limit := defaultLimit(int(req.Limit))
	offset := int(req.Offset)

	_, err := s.Sqlc.SearchArtists(ctx, sqlcdb.SearchArtistsParams{
		Artist: sql.NullString{Valid: true, String: q},
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "failed to search artists")
	}
	return &v1.SearchArtistResponse{
		Artists: nil,
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

	file, err := os.Open(filepath.Join(s.OsuDir, "Songs", string(filename)))
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

	http.ServeFile(w, r, filepath.Join(s.OsuDir, "Songs", string(filename)))
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
