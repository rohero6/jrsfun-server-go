package api

import (
	"jrsfun-server-go/model"
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

func Movie(c *gin.Context) {
	movieList := []model.MovieSite{
		{Icon: "https://ddys.pro/favicon-32x32.png", Url: "https://ddys.pro/"},
		{Icon: "https://xiaoxiaojia.oss-cn-shanghai.aliyuncs.com/statics/img/logo2.png", Url: "https://www.libvio.me/"},
		{Icon: "https://zxzjbackup.oss-cn-shenzhen.aliyuncs.com/logo.png", Url: "https://www.zxzjhd.com/"},
	}
	c.JSON(http.StatusOK, gin.H{"movie_sites": movieList})
}
