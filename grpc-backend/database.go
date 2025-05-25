package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/juli0n21/go-osu-parser/parser"
	_ "modernc.org/sqlite"
)

var ErrBeatmapCountNotMatch = errors.New("beatmap count not matching")

var osuDB *parser.OsuDB
var osuRoot string

func initDB(connectionString string, osuDb *parser.OsuDB, osuroot string) (*sql.DB, error) {

	osuDB = osuDb
	osuRoot = osuroot

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

	_, err = db.Exec("PRAGMA temp_store = MEMORY;")
	if err != nil {
		log.Fatal(err)
	}

	if err = createDB(db); err != nil {
		return nil, err
	}

	if err = checkhealth(db, osuDB); err != nil {
		if err = rebuildBeatmapDb(db, osuDB); err != nil {
			return nil, err
		}
	}

	if err = createCollectionDB(db); err != nil {
		return nil, err
	}

	collectionDB, err := parser.ParseCollectionsDB(path.Join(osuRoot, "collection.db"))
	if err != nil {
		return nil, err
	}

	if err = checkCollectionHealth(db, collectionDB); err != nil {
		if err = rebuildCollectionDb(db, collectionDB); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func createDB(db *sql.DB) error {

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS Beatmap (
    BeatmapId INTEGER DEFAULT 0,
    Artist TEXT DEFAULT '?????',
    ArtistUnicode TEXT DEFAULT '?????',
    Title TEXT DEFAULT '???????',
    TitleUnicode TEXT DEFAULT '???????',
    Creator TEXT DEFAULT '?????',
    Difficulty TEXT DEFAULT '1',
    Audio TEXT DEFAULT 'unknown.mp3',
    MD5Hash TEXT DEFAULT '00000000000000000000000000000000',
    File TEXT DEFAULT 'unknown.osu',
    RankedStatus TEXT DEFAULT Unknown,
    LastModifiedTime DATETIME DEFAULT '0001-01-01 00:00:00',
    TotalTime INTEGER DEFAULT 0,
    AudioPreviewTime INTEGER DEFAULT 0,
    BeatmapSetId INTEGER DEFAULT -1,
    Source TEXT DEFAULT '',
    Tags TEXT DEFAULT '',
    LastPlayed DATETIME DEFAULT '0001-01-01 00:00:00',
    Folder TEXT DEFAULT 'Unknown Folder',
    UNIQUE (Artist, Title, MD5Hash, Difficulty)
	);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_beatmap_md5hash ON Beatmap(MD5Hash);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_beatmap_lastModifiedTime ON Beatmap(LastModifiedTime);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_beatmap_title_artist ON Beatmap(Title, Artist);")
	if err != nil {
		return err
	}

	return nil
}

func createCollectionDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Collection (
		Name TEXT DEFAULT '',
		MD5Hash TEXT DEFAULT '00000000000000000000000000000000'
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_collection_name ON Collection(Name);")
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_collection_md5hash ON Collection(MD5Hash);")
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

func rebuildBeatmapDb(db *sql.DB, osuDb *parser.OsuDB) error {

	if _, err := db.Exec("DROP TABLE Beatmap"); err != nil {
		return err
	}

	if err := createDB(db); err != nil {
		return err
	}
	stmt, err := db.Prepare(`
		INSERT INTO Beatmap (
			BeatmapId, Artist, ArtistUnicode, Title, TitleUnicode, Creator, 
			Difficulty, Audio, MD5Hash, File, RankedStatus, 
			LastModifiedTime, TotalTime, AudioPreviewTime, BeatmapSetId, 
			Source, Tags, LastPlayed, Folder
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

func checkCollectionHealth(db *sql.DB, collectionDB *parser.Collections) error {
	rows, err := db.Query(`SELECT COUNT(*) FROM Collection GROUP BY Name;`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	if err = rows.Scan(&count); err != nil {
		return err
	}

	if count != int(collectionDB.NumberOfCollections) {
		return errors.New("Collection Count Not Matching")
	}

	rows, err = db.Query(`SELECT COUNT(*) FROM Collection;`)
	if err != nil {
		return err
	}

	if err = rows.Scan(&count); err != nil {
		return err
	}

	sum := 0
	for _, col := range collectionDB.Collections {
		sum += len(col.Beatmaps)
	}

	if count != int(sum) {
		return errors.New("Beatmap count missmatch rebuilding collections")
	}

	return nil
}

func rebuildCollectionDb(db *sql.DB, collectionDb *parser.Collections) error {
	if _, err := db.Exec("DROP TABLE Collection"); err != nil {
		return err
	}

	if err := createCollectionDB(db); err != nil {
		return err
	}

	stmt, err := db.Prepare(`
		INSERT INTO Collection (
			Name,
			MD5Hash
		) VALUES (?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt = tx.Stmt(stmt)

	for _, col := range collectionDb.Collections {
		for _, hash := range col.Beatmaps {
			_, err := stmt.Exec(col.Name, hash)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return tx.Commit()
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
	rows, err := db.Query("SELECT BeatmapId, MD5Hash, Title, Artist, Creator, Folder, File, Audio, TotalTime FROM Beatmap GROUP BY Folder ORDER BY LastModifiedTime DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return []Song{}, err
	}
	defer rows.Close()

	return scanSongs(rows)
}

func getSearch(db *sql.DB, q string, limit, offset int) (ActiveSearch, error) {
	rows, err := db.Query(
		`
	SELECT 
		BeatmapId, 
		MD5Hash, 
		Title, 
		Artist, 
		Creator, 
		Folder, 
		File, 
		Audio, 
		TotalTime 
	FROM Beatmap WHERE Title LIKE ? OR Artist LIKE ? LIMIT ? OFFSET ?
	`, "%"+q+"%", "%"+q+"%", limit, offset)
	if err != nil {
		return ActiveSearch{}, err
	}
	defer rows.Close()
	s, err := scanSongs(rows)
	if err != nil {
		return ActiveSearch{}, err
	}
	return ActiveSearch{Songs: s}, nil
}

func getArtists(db *sql.DB, q string, limit, offset int) ([]Artist, error) {
	rows, err := db.Query("SELECT Artist, COUNT(Artist) FROM Beatmap WHERE Artist LIKE ? OR Title LIKE ? GROUP BY Artist LIMIT ? OFFSET ?", "%"+q+"%", "%"+q+"%", limit, offset)
	if err != nil {
		return []Artist{}, err
	}
	defer rows.Close()

	artist := []Artist{}
	for rows.Next() {
		var a string
		var c int
		err := rows.Scan(&a, &c)
		if err != nil {
			return []Artist{}, err
		}
		artist = append(artist, Artist{Artist: a, Count: c})
	}

	return artist, nil
}

func getFavorites(db *sql.DB, q string, limit, offset int) ([]Song, error) {
	rows, err := db.Query("SELECT * FROM Songs WHERE IsFavorite = 1 LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSongs(rows)
}

func getCollection(db *sql.DB, limit, offset, index int) (Collection, error) {
	rows, err := db.Query(`
	WITH cols AS (
	SELECT 
		c.Name, 
		ROW_NUMBER() OVER (ORDER BY c.Name) AS RowNumber
	FROM Collection c
	GROUP BY c.Name
	)

	SELECT 
		c.Name, b.BeatmapId, b.MD5Hash, b.Title, b.Artist, 
		b.Creator, b.Folder, b.File, b.Audio, b.TotalTime
		FROM Collection c
		Join Beatmap b ON c.MD5Hash = b.MD5Hash
		WHERE c.Name = (SELECT Name FROM cols WHERE RowNumber = ?)
		LIMIT ? 
		OFFSET ?;`, index, limit, offset)
	if err != nil {
		return Collection{}, err
	}
	defer rows.Close()

	var c Collection
	for rows.Next() {
		s := Song{}
		if err := rows.Scan(&c.Name, &s.BeatmapID, &s.MD5Hash, &s.Title, &s.Artist, &s.Creator, &s.Folder, &s.File, &s.Audio, &s.TotalTime); err != nil {
			return Collection{}, err
		}

		s.Image = extractImageFromFile(osuRoot, s.Folder, s.File)

		c.Songs = append(c.Songs, s)
	}

	row := db.QueryRow(`SELECT COUNT(*) FROM Collection WHERE Name = ?`, c.Name)
	var count string
	row.Scan(&count)

	if i, err := strconv.Atoi(count); err == nil {
		c.Items = i
	}

	return c, nil
}

func getCollectionByName(db *sql.DB, limit, offset int, name string) (Collection, error) {
	rows, err := db.Query(`
	SELECT 
		c.Name, b.BeatmapId, b.MD5Hash, b.Title, b.Artist, 
		b.Creator, b.Folder, b.File, b.Audio, b.TotalTime
		FROM Collection c
		Join Beatmap b ON c.MD5Hash = b.MD5Hash
		WHERE c.Name = ?
		LIMIT ? 
		OFFSET ?;`, name, limit, offset)
	if err != nil {
		return Collection{}, err
	}
	defer rows.Close()

	var c Collection
	for rows.Next() {
		s := Song{}
		if err := rows.Scan(&c.Name, &s.BeatmapID, &s.MD5Hash, &s.Title, &s.Artist, &s.Creator, &s.Folder, &s.File, &s.Audio, &s.TotalTime); err != nil {
			return Collection{}, err
		}

		s.Image = extractImageFromFile(osuRoot, s.Folder, s.File)

		c.Songs = append(c.Songs, s)
	}

	row := db.QueryRow(`SELECT COUNT(*) FROM Collection WHERE Name = ?`, c.Name)
	var count string
	row.Scan(&count)

	if i, err := strconv.Atoi(count); err == nil {
		c.Items = i
	} else {
		return Collection{}, err
	}

	return c, nil
}

func getCollections(db *sql.DB, q string, limit, offset int) ([]CollectionPreview, error) {
	rows, err := db.Query(`
		SELECT 
		c.Name, 
		COUNT(b.MD5Hash) AS Count,
		MIN(b.Folder) AS Folder,
		MIN(b.File) AS File
		FROM Collection c
		JOIN Beatmap b ON c.MD5Hash = b.MD5Hash
		WHERE c.Name LIKE ?
		GROUP BY c.Name
		LIMIT ?
		OFFSET ?`, "%"+q+"%", limit, offset)
	if err != nil {
		return []CollectionPreview{}, err
	}
	defer rows.Close()

	var collections []CollectionPreview
	for rows.Next() {
		var c CollectionPreview
		var folder, file, count string
		if err := rows.Scan(&c.Name, &count, &folder, &file); err != nil {
			return []CollectionPreview{}, err
		}

		if i, err := strconv.Atoi(count); err == nil {
			c.Items = i
		}
		c.Image = extractImageFromFile(osuRoot, folder, file)

		collections = append(collections, c)
	}

	return collections, nil
}

func getSong(db *sql.DB, hash string) (Song, error) {

	row := db.QueryRow("SELECT BeatmapId, MD5Hash, Title, Artist, Creator, Folder, File, Audio, TotalTime FROM Beatmap WHERE MD5Hash = ?", hash)
	s, err := scanSong(row)
	return s, err

}

func scanSongs(rows *sql.Rows) ([]Song, error) {
	songs := []Song{}
	for rows.Next() {
		var s Song
		if err := rows.Scan(&s.BeatmapID, &s.MD5Hash, &s.Title, &s.Artist, &s.Creator, &s.Folder, &s.File, &s.Audio, &s.TotalTime); err != nil {
			return []Song{}, err
		}

		bm, err := parser.ParseOsuFile(fmt.Sprintf("%sSongs/%s/%s", osuRoot, s.Folder, s.File))
		if err != nil {
			fmt.Println(err)
			s.Image = fmt.Sprintf("404.png")
		} else {
			if len(bm.Events) > 1 && len(bm.Events[0].EventParams) > 1 {
				s.Image = fmt.Sprintf("%s/%s", s.Folder, strings.Trim(bm.Events[0].EventParams[0], "\""))
			}
		}

		songs = append(songs, s)
	}
	return songs, nil
}

func scanSong(row *sql.Row) (Song, error) {

	s := Song{}
	if err := row.Scan(&s.BeatmapID, &s.MD5Hash, &s.Title, &s.Artist, &s.Creator, &s.Folder, &s.File, &s.Audio, &s.TotalTime); err != nil {
		return Song{}, err
	}

	s.Image = extractImageFromFile(osuRoot, s.Folder, s.File)

	return s, nil
}

func extractImageFromFile(osuRoot, folder, file string) string {
	bm, err := parser.ParseOsuFile(filepath.Join(osuRoot, "Songs", folder, file))
	if err != nil {
		fmt.Println(err)
		return "404.png"
	}

	if bm.Version > 3 {
		if len(bm.Events) > 0 && len(bm.Events[0].EventParams) > 0 {
			return fmt.Sprintf(`%s/%s`, folder, strings.Trim(bm.Events[0].EventParams[0], "\""))
		}
	} else {
		fmt.Println(bm.Events)
	}

	return "404.png"
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
