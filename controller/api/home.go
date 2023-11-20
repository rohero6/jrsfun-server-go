package api

import (
	"jrsfun-server-go/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	data := service.HttpGetHomeData()
	c.JSON(http.StatusOK, gin.H{"items": data})
}

func Live(c *gin.Context) {
	link := c.Query("id")
	if link == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid 'id' parameter"})
		return
	}
	data, err := service.HttpGetLiveData(link)
	if err != nil {
		log.Printf("get live data error: %v", err)
	}
	c.JSON(http.StatusOK, data)
}

func Channel(c *gin.Context) {
	link := c.Query("id")
	if link == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid 'id' parameter"})
		return
	}
	data := service.HandlerChannel(link)
	c.JSON(http.StatusOK, data)
}
