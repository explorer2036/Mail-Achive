package api

import (
	"Mail-Achive/pkg/config"
	"Mail-Achive/pkg/storage/es"
	"Mail-Achive/pkg/storage/firebase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
)

// Server - http server
type Server struct {
	settings        *config.Config
	router          *gin.Engine
	elasticHandler  *es.Handler
	firebaseHandler *firebase.Handler
}

// NewServer returns a new http server
func NewServer(settings *config.Config) *Server {
	s := &Server{
		settings: settings,
		router:   gin.New(),
	}

	// init the elastic handler
	s.elasticHandler = es.NewHandler(settings)

	// init the firebase handler
	s.firebaseHandler = firebase.NewHandler(settings)

	// init the http routes
	s.initRoutes()

	return s
}

// Start - setup the http server
func (s *Server) Start() {
	go func() {
		if err := s.router.Run(s.settings.Server.ListenAddr); err != nil {
			log.Errorf("run the server: %v", err)
		}
	}()
}

// Close - release the used resource
func (s *Server) Close() {
	if s.firebaseHandler != nil {
		s.firebaseHandler.Close()
	}
	if s.elasticHandler != nil {
		s.elasticHandler.Close()
	}
}

// init the routes for the http server
func (s *Server) initRoutes() {
	// http api: login
	s.router.POST("/login", s.login)

	// http api: health
	s.router.GET("/health", s.health)

	// http api: refresh
	s.router.GET("/refresh", s.auth(), s.refresh)

	// http api: search
	s.router.POST("/search", s.auth(), s.search)

	// http api: upload file
	s.router.POST("/upload", s.auth(), s.upload)
}

// send the response back to clients
func send(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{
		"message": msg,
	})
}

// health for the http server
func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
