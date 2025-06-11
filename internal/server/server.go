package server

import (
	gorm_db "db_blueprints/blueprints/gorm"
	"db_blueprints/internal/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	httpProduct "db_blueprints/internal/domain/product/controller/http"
	httpUser "db_blueprints/internal/domain/user/controller/http"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	db     gorm_db.IDatabase
}

func NewServer(db gorm_db.IDatabase, cfg *config.Config) *Server {
	return &Server{
		engine: gin.Default(),
		cfg:    cfg,
		db:     db,
	}
}

func (s Server) Run() error {
	if err := s.MapRoutes(); err != nil {
		log.Printf("MapRoutes Error: %v", err)
	}

	if err := s.engine.Run(fmt.Sprintf(":%s", s.cfg.HTTP_PORT)); err != nil {
		log.Printf("Running HTTP server: %v", err)
	}

	return nil
}

func (s Server) MapRoutes() error {
	routesV1 := s.engine.Group("/api")

	httpProduct.Routes(routesV1, s.db)
	httpUser.Routes(routesV1, s.db)
	return nil
}
