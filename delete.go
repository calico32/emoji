package main

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

func delete(c *gin.Context) {
	file := c.Param("file")
	if file == "" {
		c.AbortWithStatus(400)
		return
	}

	target := path.Join(dataDir, file)

	// check if file exists
	if _, err := os.Stat(target); err != nil {
		c.AbortWithStatus(404)
		return
	}

	err := os.Remove(target)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}
