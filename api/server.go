package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db2 "github.com/inudev5/go-bank/db/sqlc"
)

type Server struct {
	store  db2.Store
	router *gin.Engine
}

func NewServer(store db2.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/transfers", server.createTransfer)
	server.router = router
	return server
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
