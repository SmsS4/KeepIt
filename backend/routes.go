package main

import (
	"github.com/gin-gonic/gin"
)

func createRoutes() {

	r := gin.Default()

	r.POST("/auth/login/", func(c *gin.Context) {

	})

	r.POST("/auth/register/", func(c *gin.Context) {

	})

	r.POST("/notes/new", func(c *gin.Context) {

	})

	r.GET("/notes", func(c *gin.Context) {

	})

	r.PUT("/notes", func(c *gin.Context) {

	})

	r.DELETE("/notes", func(c *gin.Context) {

	})

	r.Run()
}
