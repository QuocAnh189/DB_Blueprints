package server

import (
	"db_blueprints/config"
	db "db_blueprints/db_sql/database"
	"fmt"
	"log"

	httpProduct "db_blueprints/db_sql/internal/domain/product/controller/http"
	httpUser "db_blueprints/db_sql/internal/domain/user/controller/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	db     db.DBTX
}

func NewServer(db db.DBTX, cfg *config.Config) *Server {
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
