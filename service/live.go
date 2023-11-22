package service

import (
	"fmt"
	"io/ioutil"
	"jrsfun-server-go/cache"
	"jrsfun-server-go/model"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func HttpGetLiveData(link string) (model.LiveResponse, error) {
	key := fmt.Sprintf(model.LiveKey, link)
	val, exist := cache.GoCache.Get(key)
	if !exist {
		data, _ := GetLiveData(link)
		if data.Streams != nil {
			go cache.GoCache.Set(key, data, 3*time.Hour)
		}
		return data, nil
	} else {
		data, _ := val.(model.LiveResponse)
		return data, nil
	}
}
func GetLiveData(link string) (model.LiveResponse, error) {
	domain := strings.Split(link, ".com")[0] + ".com"
	html := fetch(link)
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return model.LiveResponse{}, err
	}
	subBox := doc.Find(".sub_list").Find("ul")
	tempLab, _ := getLab(subBox.First())
	team := getTeam(subBox.First())

	data := model.LiveResponse{
		Lab:     tempLab,
		Streams: []model.StreamProp{},
		Team:    team,
	}
	streams := make([]model.StreamProp, 0)
	subPlayList := doc.Find("#loc-tool-player .sub_box .sub_playlist .sub_channel").Find("a")
	// 使用 WaitGroup 来同步 goroutine
	var wg sync.WaitGroup
	var mu sync.Mutex // 用于保护共享数据 streams

	// 这里使用 channel 传递数据，避免 goroutine 间的竞争
	streamCh := make(chan model.StreamProp, 10) // 可根据需求调整缓冲大小

	subPlayList.Each(func(i int, element *goquery.Selection) {
		wg.Add(1)
		go func(element *goquery.Selection) {
			defer wg.Done()

			name := strings.TrimSpace(element.Find("strong").Text())
			dataPlay, _ := element.Attr("data-play")
			dataURLID, _ := element.Attr("data-urlid")
			tempLink := domain + dataPlay + dataURLID
			if !strings.HasSuffix(dataPlay, "=") {
				tempLink = domain + dataPlay
			}
			// 将数据发送到 channel
			streamCh <- model.StreamProp{
				M3U8URL: tempLink,
				Name:    name,
			}
		}(element)
	})

	// 关闭 channel，以便下面的 range 循环结束
	go func() {
		wg.Wait()
		close(streamCh)
	}()

	// 从 channel 中读取数据并收集到 streams 切片中
	for stream := range streamCh {
		mu.Lock()
		streams = append(streams, stream)
		mu.Unlock()
	}
	go cache.GoCache.Set(fmt.Sprintf(model.StreamKey, link), streams, 5*time.Hour)
	m3u8, err := GetFirstM3u8(streams[0].M3U8URL, domain)
	if err != nil {
		return model.LiveResponse{}, err
	}
	data.Streams = []model.StreamProp{
		{Name: streams[0].Name,
			M3U8URL: m3u8},
	}
	return data, nil
}

func fetch(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func GetFirstM3u8(rawUrl string, domain string) (string, error) {
	id := strings.Split(strings.Split(rawUrl, "id=")[1], "&")[0]
	playUrl := fmt.Sprintf("%s/play/%s.html", domain, id)

	resp, err := http.Get(playUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GET request failed: %s", resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	var m3u8 string
	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			if strings.Contains(src, "id=") {
				m3u8 = strings.Split(src, "id=")[1]
			}
		}
	})
	return m3u8, nil
}
