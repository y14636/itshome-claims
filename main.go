package main

import (
    "path"
    "path/filepath"

    "github.com/gin-gonic/gin"

    "github.com/y14636/itshome-claims/handlers"
)

func main() {
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
    r.POST("/claims", handlers.AddClaimsHandler)
    r.POST("/modifiedclaims", handlers.AddModifiedClaimsHandler)
    r.DELETE("/modifiedclaims/:id", handlers.DeleteClaimsHandler)

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