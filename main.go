package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	loadEnvironmentVariables()

	store, err := PostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.InitializeStorage(); err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	serverStr := "localhost:" + port
	server := Server(serverStr, store)
	server.Run()
}

func loadEnvironmentVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error while loading .env file")
	}
}
