package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imloama/syncrepo/sync"
)

func main() {
	sync.DumpConfig()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/repos", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": sync.GetConfig().Repos,
		})
	})
	r.GET("/sync", func(c *gin.Context) {
		repo := c.Query("repo")
		err := sync.GitService.Sync(repo)
		if err!=nil{
			c.JSON(400, gin.H{ "message": err.Error() })
		}
		c.JSON(200, gin.H{ "message": "ok" })
	})

	cfg := sync.GetConfig()
	ip := cfg.Ip
	port := cfg.Port
	if ip == "" {
		ip = "0.0.0.0"
	}
	if port <= 0 {
		port = 3773
	}
	listenon := fmt.Sprintf("%s:%d", ip, port)
	fmt.Printf("listen on %s", listenon)
	r.Run(listenon)
}
