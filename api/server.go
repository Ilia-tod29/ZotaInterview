package api

import (
	"ZotaInterview/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config util.Config
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config) (*Server, error) {
	server := &Server{
		config: config,
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	//router.Use(CORSMiddleware())

	router.POST("/deposit", s.depositMoney)
	router.GET("/status", s.checkDepositStatus)
	router.GET("/payment-return", s.paymentReturn)

	s.router = router
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
