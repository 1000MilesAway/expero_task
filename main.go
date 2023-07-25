package main

import (
	"github.com/gin-gonic/gin"
)

var HTTPServerURL = "0.0.0.0:5000"

// type Server struct {
// 	server    *http.Server
// 	router    *gin.Engine
// 	startOnce *sync.Once
// 	closeOnce *sync.Once
// 	deadChan  chan any
// }

// func (s *Server) Start() {

// 	err := errors.New("HTTP server has been started already")
// 	s.startOnce.Do(func() {
// 		defer close(s.deadChan)

// 		log.Info(s, "Starting listening")
// 		if err = s.server.ListenAndServe(); err != nil {
// 			log.Warning(s, err.Error())
// 			err = nil
// 		}
// 	})
// 	if err != nil {
// 		log.Error(s, err.Error())
// 	}
// }

// func (s *Server) Close() {
// 	s.closeOnce.Do(func() {
// 		log.Warning(s, "Stopping and closing")
// 		s.server.Close()
// 	})
// }

// func (s *Server) Dead() <-chan any {
// 	return s.deadChan
// }

// func (s *Server) String() string {
// 	return fmt.Sprintf("HTTP SERVER addr=%s", HTTPServerURL)
// }

// var instance *Server
// var serverOnce = &sync.Once{}

// func GetServer() *Server {
// 	serverOnce.Do(func() {
// 		gin.SetMode(gin.ReleaseMode)
// 		router := gin.New()
// 		router.Use(
// 			func(c *gin.Context) {
// 				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, application/json")
// 				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
// 				if c.Request.Method == "OPTIONS" {
// 					c.AbortWithStatus(200)
// 					return
// 				}
// 				c.Next()
// 			})
// 		router.Use(gin.Recovery())
// 		router.POST("/aggregation", instance.Aggregation)

// 		instance = &Server{
// 			server: &http.Server{
// 				Addr:    HTTPServerURL,
// 				Handler: router,
// 			},
// 			router:    router,
// 			deadChan:  make(chan any),
// 			startOnce: &sync.Once{},
// 			closeOnce: &sync.Once{},
// 		}
// 	})
// 	return instance
// }

func main() {
	defer GetDB().Close()

	router := gin.Default()
	router.POST("/aggregation", Aggregation)
	router.POST("/deaggregation", DeAggregation)

	router.Run("localhost:7777")
	// GetServer().Start()
	// go GetServer().Start()
	// defer GetServer().Close()
}
