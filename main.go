package main

import (
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/y14636/itshome-claims/handlers"
	logConfig "github.com/y14636/itshome-claims/logging"
)

func main() {
	logConfig.InitializeLogging("./logs/claims.log")

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	r.GET("/claims", handlers.GetClaimsListHandler)
	r.GET("/modifiedclaims", handlers.GetModifiedClaimsListHandler)
	r.GET("/searchclaims/:search", handlers.GetClaimsResultsHandler)
	r.GET("/claims/:claimsId", handlers.GetClaimsListByIdHandler)
	//r.POST("/claims", handlers.AddClaimsHandler)
	r.GET("/modifiedclaims/:claimsData", handlers.AddClaimsHandler)
	r.DELETE("/modifiedclaims/:id", handlers.DeleteClaimsHandler)
	r.POST("/logging", handlers.LogWebMessages)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, POST, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
