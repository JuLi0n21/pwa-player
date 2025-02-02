package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/juli0n21/go-osu-parser/parser"
	_ "modernc.org/sqlite"
)

var ErrBeatmapCountNotMatch = errors.New("beatmap count not matching")

func initDB(connectionString string, osuDb *parser.OsuDB) (*sql.DB, error) {

	dir := filepath.Dir(connectionString)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	db, err := sql.Open("sqlite", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database %s: %v", connectionString, err)
	}

	if err = createDB(db); err != nil {
		return nil, err
	}

	if err = checkhealth(db, osuDb); err != nil {
		if err = rebuildDb(db, osuDb); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func createDB(db *sql.DB) error {
	_, err := db.Query(`
	CREATE TABLE IF NOT EXISTS Beatmap (
    BeatmapId INTEGER DEFAULT 0,
    Artist TEXT DEFAULT '?????',
    ArtistUnicode TEXT DEFAULT '?????',
    Title TEXT DEFAULT '???????',
    TitleUnicode TEXT DEFAULT '???????',
    Creator TEXT DEFAULT '?????',
    Difficulty TEXT DEFAULT '1',
    AudioFileName TEXT DEFAULT 'unknown.mp3',
    MD5Hash TEXT DEFAULT '00000000000000000000000000000000',
    FileName TEXT DEFAULT 'unknown.osu',
    RankedStatus TEXT DEFAULT Unknown,
    LastModifiedTime DATETIME DEFAULT '0001-01-01 00:00:00',
    TotalTime INTEGER DEFAULT 0,
    AudioPreviewTime INTEGER DEFAULT 0,
    BeatmapSetId INTEGER DEFAULT -1,
    Source TEXT DEFAULT '',
    Tags TEXT DEFAULT '',
    LastPlayed DATETIME DEFAULT '0001-01-01 00:00:00',
    FolderName TEXT DEFAULT 'Unknown Folder',
    UNIQUE (Artist, Title, MD5Hash, Difficulty)
	);
	`)

	if err != nil {
		return err
	}

	return nil
}

func checkhealth(db *sql.DB, osuDb *parser.OsuDB) error {

	rows, err := db.Query(`SELECT COUNT(*) FROM Beatmap GROUP BY BeatmapSetId;`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	if err = rows.Scan(&count); err != nil {
		return err
	}

	if count != int(osuDb.FolderCount) {
		log.Println("Folder count missmatch rebuilding db...")
		return ErrBeatmapCountNotMatch
	}

	rows, err = db.Query(`SELECT COUNT(*) FROM Beatmap;`)
	if err != nil {
		return err
	}

	if err = rows.Scan(&count); err != nil {
		return err
	}

	if count != int(osuDb.NumberOfBeatmaps) {
		log.Println("Beatmap count missmatch rebuilding db...")
		return ErrBeatmapCountNotMatch
	}

	return nil
}

func rebuildDb(db *sql.DB, osuDb *parser.OsuDB) error {

	if _, err := db.Query("DROP TABLE Beatmap"); err != nil {
		return err
	}

	if err := createDB(db); err != nil {
		return err
	}
	stmt, err := db.Prepare(`
		INSERT INTO Beatmap (
			BeatmapId, Artist, ArtistUnicode, Title, TitleUnicode, Creator, 
			Difficulty, AudioFileName, MD5Hash, FileName, RankedStatus, 
			LastModifiedTime, TotalTime, AudioPreviewTime, BeatmapSetId, 
			Source, Tags, LastPlayed, FolderName
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	//	ON CONFLICT (Artist, Title, MD5Hash) DO NOTHING

	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt = tx.Stmt(stmt)

	for i, beatmap := range osuDb.Beatmaps {
		//fmt.Println(i, beatmap.Artist, beatmap.SongTitle, beatmap.MD5Hash)
		_, err := stmt.Exec(
			beatmap.DifficultyID, beatmap.Artist, beatmap.ArtistUnicode,
			beatmap.SongTitle, beatmap.SongTitleUnicode, beatmap.Creator,
			beatmap.Difficulty, beatmap.AudioFileName, beatmap.MD5Hash,
			beatmap.FileName, beatmap.RankedStatus, beatmap.LastModificationTime,
			beatmap.TotalTime, beatmap.AudioPreviewStartTime, beatmap.BeatmapID,
			beatmap.SongSource, beatmap.SongTags, beatmap.LastPlayed, beatmap.FolderName,
		)
		if err != nil {
			fmt.Println(i, "hash: ", beatmap.MD5Hash, "artist:", beatmap.Artist, "title:", beatmap.SongTitle, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func getBeatmapCount(db *sql.DB) int {
	rows, err := db.Query("SELECT COUNT(*) FROM Beatmap")
	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println(err)
			return 0
		}
	} else {
		return 0
	}

	return count
}

func getRecent(db *sql.DB, limit, offset int) ([]Song, error) {
	rows, err := db.Query("SELECT * FROM Songs ORDER BY LastPlayed DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSongs(rows)
}

func getSearch(db *sql.DB, q string, limit, offset int) (ActiveSearch, error) {
	rows, err := db.Query("SELECT * FROM Songs WHERE Title LIKE ? OR Artist LIKE ? LIMIT ? OFFSET ?", "%"+q+"%", "%"+q+"%", limit, offset)
	if err != nil {
		return ActiveSearch{}, err
	}
	defer rows.Close()
	_, err = scanSongs(rows)
	if err != nil {
		return ActiveSearch{}, err
	}
	return ActiveSearch{}, nil
}

func getArtists(db *sql.DB, q string, limit, offset int) ([]string, error) {
	rows, err := db.Query("SELECT * FROM Songs WHERE Title LIKE ? OR Artist LIKE ? LIMIT ? OFFSET ?", "%"+q+"%", "%"+q+"%", limit, offset)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	_, err = scanSongs(rows)
	if err != nil {
		return []string{}, err
	}
	return []string{}, nil
}

func getFavorites(db *sql.DB, q string, limit, offset int) ([]Song, error) {
	rows, err := db.Query("SELECT * FROM Songs WHERE IsFavorite = 1 LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSongs(rows)
}

func getCollection(db *sql.DB, index int) (Collection, error) {
	row := db.QueryRow("SELECT * FROM Collections WHERE CollectionId = ?", index)
	return scanCollection(row)
}

func getCollections(db *sql.DB, q string, limit, offset int) ([]Collection, error) {
	//not correct
	rows, err := db.Query("SELECT * FROM Collections WHERE name = ? LIMIT ? OFFSET ?", q, limit, offset)
	if err != nil {
		return []Collection{}, err
	}
	return scanCollections(rows)
}

func getSong(db *sql.DB, hash string) (Song, error) {
	row := db.QueryRow("SELECT * FROM Songs WHERE MD5Hash = ?", hash)
	return scanSong(row)
}

func scanSongs(rows *sql.Rows) ([]Song, error) {

	var songs []Song
	for rows.Next() {
		var s Song
		if err := rows.Scan(&s); err != nil {
			return []Song{}, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

func scanSong(row *sql.Row) (Song, error) {

	var s Song
	if err := row.Scan(&s); err != nil {
		return Song{}, err
	}
	return s, nil
}

func scanCollections(rows *sql.Rows) ([]Collection, error) {

	var collection []Collection
	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c); err != nil {
			return []Collection{}, err
		}
		collection = append(collection, c)
	}
	return collection, nil
}

func scanCollection(row *sql.Row) (Collection, error) {

	var c Collection
	if err := row.Scan(&c); err != nil {
		return Collection{}, err
	}
	return c, nil
}

func scanCollectionPreviews(rows *sql.Rows) ([]CollectionPreview, error) {

	var collection []CollectionPreview
	for rows.Next() {
		var c CollectionPreview
		if err := rows.Scan(&c); err != nil {
			return []CollectionPreview{}, err
		}
		collection = append(collection, c)
	}
	return collection, nil
}

func scanCollectionPreview(row *sql.Row) (CollectionPreview, error) {

	var c CollectionPreview
	if err := row.Scan(&c); err != nil {
		return CollectionPreview{}, err
	}
	return c, nil
}
