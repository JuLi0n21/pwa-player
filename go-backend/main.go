package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/juli0n21/go-osu-parser/parser"
)

func main() {

	filename := "/mnt/g/Anwendungen/osu!/osu!.db"

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	osuDb, err := parser.ParseOsuDB(filename)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB("./data/music.db", osuDb)
	if err != nil {
		log.Fatal(err)
	}

	run(db, osuDb)
}
