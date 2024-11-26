package main

import (
	"log"
	"t1/internal/pkg/db"
	"t1/internal/pkg/server"
)

func main() {
	// get port from environment
	s := server.New(":8090")

	err := db.Migrate(s.DB, "scripts/migration.sql")
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	s.Start()
}
