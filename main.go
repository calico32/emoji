package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var mode string
var dataDir string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: ", err)
		log.Print("Continuing anyway!")
	}
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()

	e.TrustedPlatform = gin.PlatformCloudflare
	e.SetTrustedProxies([]string{"172.0.0.1/16"})

	dataDirVar := os.Getenv("DATA_DIR")
	if dataDirVar == "" {
		log.Println("DATA_DIR not set, using default")
		dataDirVar = "./data"
	}

	manageKey := os.Getenv("MANAGE_KEY")
	if manageKey == "" {
		log.Println("MANAGE_KEY not set!")
		os.Exit(1)
	}

	dataDir, err = filepath.Abs(path.Join(dataDirVar))
	if err != nil {
		log.Println("Error getting absolute path: ", err)
		os.Exit(1)
	}

	authenticated := func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header != "Bearer "+manageKey {
			c.AbortWithStatus(401)
			return
		}
	}

	fileServer := http.FileServer(gin.Dir(dataDir, false))

	indexRedirect := os.Getenv("INDEX_REDIRECT")

	e.POST("/", authenticated, upload)
	e.DELETE("/:file", authenticated, delete)
	e.OPTIONS("/", authenticated, list)
	e.GET("/*path", func(c *gin.Context) {
		path := c.Param("path")
		segments := []string{}
		for _, segment := range strings.Split(path, "/") {
			if segment != "" {
				if segment == ".." {
					c.AbortWithStatus(400)
					return
				}
				segments = append(segments, segment)
			}
		}

		switch len(segments) {
		case 0:
			if indexRedirect != "" {
				c.Redirect(302, indexRedirect)
			}
		case 1:
			fileServer.ServeHTTP(c.Writer, c.Request)
		case 2:
			resize(c, segments[0], segments[1])
		default:
			c.AbortWithStatus(400)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Run(":" + port)
}
