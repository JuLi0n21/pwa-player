package main

import (
	sqlcdb "backend/internal/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/juli0n21/go-osu-parser/parser"
	_ "modernc.org/sqlite"
)

var ErrBeatmapCountNotMatch = errors.New("beatmap count not matching")

var osuDB *parser.OsuDB
var osuRoot string

func initDB(connectionString string, osuDb *parser.OsuDB, osuroot string) (*sql.DB, *sqlcdb.Queries, error) {

	osuDB = osuDb
	osuRoot = osuroot

	dir := filepath.Dir(connectionString)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	db, err := sql.Open("sqlite", connectionString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open SQLite database %s: %v", connectionString, err)
	}

	_, err = db.Exec("PRAGMA temp_store = MEMORY;")
	if err != nil {
		log.Fatal(err)
	}

	if err = createDB(db); err != nil {
		return nil, nil, err
	}

	sqlcQueries := sqlcdb.New(db)

	if err = checkhealth(sqlcQueries, osuDB); err != nil {
		if err = rebuildBeatmapDb(db, osuDB); err != nil {
			return nil, nil, err
		}
	}

	if err = createCollectionDB(db); err != nil {
		return nil, nil, err
	}

	collectionDB, err := parser.ParseCollectionsDB(path.Join(osuRoot, "collection.db"))
	if err != nil {
		return nil, nil, err
	}

	if err = checkCollectionHealth(db, collectionDB); err != nil {
		if err = rebuildCollectionDb(db, collectionDB); err != nil {
			return nil, nil, err
		}
	}

	return db, sqlcQueries, nil
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
    LastModifiedTime  INTEGER DEFAULT 0,
    TotalTime INTEGER DEFAULT 0,
    AudioPreviewTime INTEGER DEFAULT 0,
    BeatmapSetId INTEGER DEFAULT -1,
    Source TEXT DEFAULT '',
    Tags TEXT DEFAULT '',
    LastPlayed INTEGER DEFAULT 0,
    Folder TEXT DEFAULT 'Unknown Folder',
    UNIQUE (Artist, Title, MD5Hash, Difficulty, Folder)
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

func checkhealth(db *sqlcdb.Queries, osuDb *parser.OsuDB) error {

	count, err := db.GetBeatmapSetCount(context.TODO())
	if err != nil {
		return err
	}

	if count != int64(osuDb.FolderCount) {
		log.Println("Folder count missmatch rebuilding db...")
		return ErrBeatmapCountNotMatch
	}

	count, err = db.GetBeatmapCount(context.TODO())
	if err != nil {
		return err
	}

	if count != int64(osuDb.NumberOfBeatmaps) {
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

func extractImageFromFile(osuRoot, folder, file string) string {
	bm, err := parser.ParseOsuFile(filepath.Join(osuRoot, "Songs", folder, file))
	if err != nil {
		fmt.Println(err)
		return "404.png"
	}

	bgImage := bm.BackgroundImage()
	if bgImage != "" {
		return filepath.Join(folder, bgImage)
	}

	return "404.png"
}
