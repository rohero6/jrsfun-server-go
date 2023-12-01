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
		{Icon: "https://toolsdar.cn/img/1680313844-Pasted.png", Url: "https://ddys.pro/", Name: "低端影视"},
		{Icon: "https://toolsdar.cn/img/1676522049-Pasted.png", Url: "https://www.libvio.me/", Name: "LIBVIO影视"},
		{Icon: "https://toolsdar.cn/img/1673323796-Pasted.png", Url: "https://www.zxzjhd.com/", Name: "在线之家"},
		{Icon: "https://api.iowen.cn/favicon/cokemv.link.png", Url: "https://cokemv.link/", Name: "COKEMV影视"},
		{Icon: "https://www.mjtt.fun/template/mojia/asset/img/logo.png", Url: "https://www.mjtt.fun/", Name: "美剧天堂"},
		{Icon: "https://dianyi.ng/mxstatic/image/logo.png", Url: "https://dianyi.ng/", Name: "电影先生"},
	}
	c.JSON(http.StatusOK, gin.H{"movie_sites": movieList})
}
