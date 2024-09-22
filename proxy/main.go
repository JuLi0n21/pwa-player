package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if value := os.Getenv("ENV"); value == "prod" {
		if ok := godotenv.Load(".env"); ok != nil {
			panic(".env not found")
		}
	} else {
		fmt.Println("Fallback to dev.env")
		if ok := godotenv.Load("dev.env"); ok != nil {
			panic("dev.env not found")
		}

	}
	InitDB()

	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
