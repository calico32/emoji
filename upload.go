package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	filename := file.Filename
	if filename == "" {
		c.AbortWithStatus(400)
		return
	}

	// remove file ext
	filename = filename[:len(filename)-len(filepath.Ext(filename))]

	target := path.Join(dataDir, filename)

	// check if file exists
	if _, err := os.Stat(target); err == nil {
		c.AbortWithStatus(409)
		return
	}

	err = c.SaveUploadedFile(file, target)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"path":   "/" + filename,
	})
}
