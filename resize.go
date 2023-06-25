package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

func resize(c *gin.Context, file, sizeStr string) {
	size, err := strconv.ParseInt(sizeStr, 10, 16)
	if err != nil {
		c.AbortWithStatus(400)
	}

	target := path.Join(dataDir, file)

	// check if file exists
	if _, err := os.Stat(target); err != nil {
		c.AbortWithStatus(404)
		return
	}

	data, err := os.ReadFile(target)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}

	image.RegisterFormat("jpeg", "\xff\xd8\xff", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "\x89PNG\x0d\x0a\x1a\x0a", png.Decode, png.DecodeConfig)
	image.RegisterFormat("webp", "RIFF????WEBP", webp.Decode, webp.DecodeConfig)
	image.RegisterFormat("gif", "GIF8?a", gif.Decode, gif.DecodeConfig)

	// resize
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}

	output := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))
	draw.CatmullRom.Scale(output, output.Bounds(), img, img.Bounds(), draw.Over, nil)

	c.Header("Content-Type", "image/png")
	err = png.Encode(c.Writer, output)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(500)
		return
	}

	c.Status(200)
}
