package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {

	if ok := godotenv.Load("dev.env"); ok != nil {
		panic(".env not found")
	}

	InitDB()

	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
