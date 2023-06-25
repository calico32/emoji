package main

import (
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func list(c *gin.Context) {
	// list files
	files, err := filepath.Glob(path.Join(dataDir, "*"))
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	names := make([]string, len(files))
	for i, file := range files {
		names[i] = path.Base(file)
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"files":  names,
	})
}
