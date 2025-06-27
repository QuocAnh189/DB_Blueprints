package main

import (
	"db_blueprints/config"
	db "db_blueprints/gorm/database"
	"db_blueprints/gorm/internal/server"
	"fmt"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {
	cfg := config.LoadConfig()

	database, err := db.NewDatabase(cfg)
	if err != nil {
		fmt.Println("Cannot connect to database", err)
	}

	httpSvr := server.NewServer(database, cfg)

	wg.Add(1)

	// Run HTTP server
	go func() {
		defer wg.Done()
		if err := httpSvr.Run(); err != nil {
			log.Println("Running HTTP server error:", err)
		}
	}()

	wg.Wait()

	println(database)
}
