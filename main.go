package main

import "log"

func main() {
	store, err := PostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := Server("localhost:8088", store)
	server.Run()
}
