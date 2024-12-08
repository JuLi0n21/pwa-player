package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if value := os.Getenv("ENV"); value == "prod" {
		if ok := godotenv.Load(".env"); ok != nil {
			log.Println(".env not found")
		}
	} else {
		fmt.Println("Fallback to dev.env")
		if ok := godotenv.Load("dev.env"); ok != nil {
			log.Println("dev.env not found, falling back to ENVIORMENT VARS")
		}

	}
	InitDB()

	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
