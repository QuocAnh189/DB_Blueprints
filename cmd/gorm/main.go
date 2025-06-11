package main

import (
	"db_blueprints/internal/config"
	"fmt"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)
}
