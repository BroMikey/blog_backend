package api

import (
	db "github.com/BroMikey/blog_backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()
	server.router = router

	// add routes here
	router.POST("/coin", server.createCoin)
	router.GET("/coin/:id", server.getCoin)
	router.GET("/coin", server.listCoin)

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
