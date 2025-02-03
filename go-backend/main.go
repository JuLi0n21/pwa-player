package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/juli0n21/go-osu-parser/parser"
)

// @title go-osu-music-hoster
// @version 1.0
// @description Server Hosting ur own osu files over a simple Api

// @host localhost:8080
// @BasePath /api/v1/
func main() {

	filename := "/mnt/g/Anwendungen/osu!/osu!.db"
	osuRoot := "/mnt/g/Anwendungen/osu!/"

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	osuDb, err := parser.ParseOsuDB(filename)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB("./data/music.db", osuDb, osuRoot)
	if err != nil {
		log.Fatal(err)
	}

	s := &Server{
		Port:   ":8080",
		Db:     db,
		OsuDb:  osuDb,
		OsuDir: osuRoot,
	}

	run(s)
}
