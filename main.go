package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println(42)
}
