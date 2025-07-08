package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yosa/ocr-golang-back/db"
	"github.com/yosa/ocr-golang-back/token"
	"github.com/yosa/ocr-golang-back/util"
)

type Server struct {
	queries    *db.Queries
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, queries *db.Queries) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		queries:    queries,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	router.POST("/users", server.CreateUserHandler)
	router.POST("/users/login", server.LoginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	server.router = router
	return server, nil

}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
