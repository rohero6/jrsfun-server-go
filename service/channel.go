package service

import (
	"fmt"
	"jrsfun-server-go/cache"
	"jrsfun-server-go/manager"
	"jrsfun-server-go/model"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

func HandlerChannel(link string) model.ChannelResp {
	key := fmt.Sprintf(model.ChannelKey, link)
	domain := strings.Split(link, ".com")[0] + ".com"
	val, exist := cache.GoCache.Get(key)
	if !exist {
		streamKey := fmt.Sprintf(model.StreamKey, link)
		steamData, found := cache.GoCache.Get(streamKey)
		if found {
			steamProp, _ := steamData.([]model.StreamProp)
			if len(steamProp) > 0 {
				data := GetChannel(&steamProp, domain)
				var res model.ChannelResp
				for _, val := range data.Streams {
					if val.M3U8URL != "" {
						res.Streams = append(res.Streams, val)
					}
				}
				// 对名字进行排序，带有"主播"的排在后面
				sort.Slice(res.Streams, func(i, j int) bool {
					return strings.Contains(res.Streams[i].Name, "主播") && !strings.Contains(res.Streams[j].Name, "主播")
				})
				if res.Streams != nil {
					go cache.GoCache.Set(key, res, 60*time.Minute)
				}
				return res
			}
		}

	} else {
		data, _ := val.(model.ChannelResp)
		return data
	}
	return model.ChannelResp{}
}
func GetChannel(streams *[]model.StreamProp, domain string) model.ChannelResp {
	var wg sync.WaitGroup
	wg.Add(len(*streams))
	for i := range *streams {
		go func(idx int) {
			defer wg.Done()
			page, _ := manager.Context.NewPage()
			m3u8 := NavigateAndHandleLogic(page, (*streams)[idx].M3U8URL, domain)
			(*streams)[idx].M3U8URL = m3u8
			if page != nil {
				err := page.Close()
				if err != nil {
					log.Printf("err get channel")
				}
			}

		}(i)
	}
	wg.Wait()
	return model.ChannelResp{Streams: *streams}
}
func NavigateAndHandleLogic(page playwright.Page, url string, domain string) string {
	page.Goto(url)
	html, _ := page.Content()
	reader := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(reader)
	iframes := doc.Find("#myElement iframe").First()
	if iframes != nil {
		src, exists := iframes.Attr("src")
		if exists && (!strings.Contains(src, "m3u8")) {
			page.Goto(domain + src)
			html, _ := page.Content()
			reader := strings.NewReader(html)
			doc, _ := goquery.NewDocumentFromReader(reader)
			src = doc.Find("iframe").First().AttrOr("src", "")
		}
		fmt.Println(src)
		var m3u8 string
		split := strings.Split(src, "id=")
		if len(split) > 1 {
			m3u8 = split[1]
		}
		if len(m3u8) > 0 && !strings.HasPrefix(m3u8, "http") {
			m3u8 = "https:" + m3u8
		}
		return m3u8
	}
	return ""
}
