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
		{Icon: "https://yingshi.tv/static/images/logo.svg", Url: "https://yingshi.tv/", Name: "影视TV"},
		{Icon: "https://hoho.tv/image/logo.svg", Url: "https://hoho.tv/", Name: "HOHOTV"},
		{Icon: "https://www.wjys.cc/mxstatic/image/index_logo.png", Url: "https://www.wjys.cc/", Name: "万佳影视"},
		{Icon: "https://www.physkan.com/logo/bphys.png", Url: "https://www.physkan.com/", Name: "胖虎影视"},
		{Icon: "https://edu-image.nosdn.127.net/E4789315155514D7814B2420CDBC4DDF.png", Url: "https://www.6080yy3.com/", Name: "新视觉影院"},
		{Icon: "https://img1.baidu.com/it/u=4279398509,1916506432&fm=253&app=138&size=w931&n=0&f=JPEG&fmt=auto?sec=1702659600&t=925a866a434fb1c2940a98eda6ac16bc", Url: "https://www.missav.com/", Name: "同花顺"},
	}
	c.JSON(http.StatusOK, gin.H{"movie_sites": movieList})
}
