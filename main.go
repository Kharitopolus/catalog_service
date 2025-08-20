package main

import "log"

func main() {
	store, err := NewPostgresStore(
		"postgres://mb_catalog_user:hackme@0.0.0.0:5432/mb_catalog_db?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	NewAPIServer(":8000", store).Run()
}
