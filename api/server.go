package api

import (
	db "github.com/BroMikey/blog_backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {

	servicer := &Server{store: store}
	router := gin.Default()
	servicer.router = router

	// add routes here
	router.POST("/")

	return servicer
}
