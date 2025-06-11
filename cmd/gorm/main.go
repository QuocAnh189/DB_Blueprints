package main

import (
	gorm_db "db_blueprints/blueprints/gorm"
	"db_blueprints/internal/config"
	"db_blueprints/internal/server"
	"fmt"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {
	cfg := config.LoadConfig()

	database, err := gorm_db.NewDatabase(cfg)
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
