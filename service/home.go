package service

import (
	"fmt"
	"jrsfun-server-go/cache"
	"jrsfun-server-go/manager"
	"jrsfun-server-go/model"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func HttpGetHomeData() []model.Item {
	key := "jrsfun_home"
	val, exist := cache.GoCache.Get(key)
	if !exist {
		data := GetHomeData()
		if data != nil {
			go cache.GoCache.Set(key, data, 5*time.Minute)
		}
		return data
	} else {
		data, _ := val.([]model.Item)
		return data
	}
}
func GetHomeData() []model.Item {
	url := "http://m.jrsbxj.com/"
	page, err := manager.Context.NewPage()
	defer page.Close()
	if err != nil {
		log.Printf("new page has err: %v", err)
		return nil
	}
	page.Goto(url)
	html, err := page.Content()
	if err != nil {
		log.Printf("page content retrieval failed: %v", err)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("goquery.NewDocumentFromReader has err: %v", err)
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	data := make([]model.Item, 0)

	subLists := doc.Find(".loc_match_list .sub_list")
	subLists.Each(func(i int, subList *goquery.Selection) {
		wg.Add(1)
		go func(subList *goquery.Selection) {
			defer wg.Done()
			hot := i == 0

			subList.Find("ul").Each(func(i int, element *goquery.Selection) {
				lab, subStatus := getLab(element)
				team := getTeam(element)
				zhibos := getZhibos(element)
				isNotStart := subStatus == "未开始"

				item := model.Item{
					Lab:        lab,
					Team:       team,
					Zhibos:     zhibos,
					IsNotStart: isNotStart,
					SubStatus:  subStatus,
				}
				item.Lab.Hot = hot

				mu.Lock()
				data = append(data, item)
				mu.Unlock()
			})
		}(subList)
	})

	wg.Wait()
	return data
}

func handleSubStatus(subStatus string) string {
	switch subStatus {
	case "中":
		return "进行中"
	case "未":
		return "未开始"
	}
	return subStatus
}

func getLab(element *goquery.Selection) (model.Lab, string) {
	labEvent := element.Find(".lab_events .name").Text()
	labTime := element.Find(".lab_time").Text()
	labBc := element.Find(".lab_bc").Text()
	labJq := element.Find(".lab_jq").Text()
	labEventStyle, _ := element.Find(".lab_events").Attr("style")
	labEventBgColor := strings.Split(labEventStyle, ":")[1]

	subStatus := handleSubStatus(element.Find(".sub_status").Text())
	dataLid, _ := element.Attr("data-lid")
	id := dataLid + "@" + labTime

	return model.Lab{
		LabEvent:        labEvent,
		LabTime:         labTime,
		LabBC:           labBc,
		LabJQ:           labJq,
		LabEventBGColor: labEventBgColor,
		BF:              "",
		ID:              id,
	}, subStatus
}

func getTeam(element *goquery.Selection) model.Team {
	teamHome := element.Find(".lab_team_home .name").Text()
	teamHomeIcon, _ := element.Find(".lab_team_home img").Attr("src")
	teamAway := element.Find(".lab_team_away .name").Text()
	teamAwayIcon, _ := element.Find(".lab_team_away img").Attr("src")
	bf := element.Find(".lab_team_home .bf").Text() + "-" + element.Find(".lab_team_away .bf").Text()

	return model.Team{
		TeamHome:     teamHome,
		TeamHomeIcon: teamHomeIcon,
		TeamAway:     teamAway,
		TeamAwayIcon: teamAwayIcon,
		Bf:           bf,
	}
}

func getZhibos(element *goquery.Selection) []model.ZhiboProp {
	zhibos := []model.ZhiboProp{}
	element.Find(".lab_channel .ok").Each(func(i int, channel *goquery.Selection) {
		i++
		name := fmt.Sprintf("线路%v", i)
		url, _ := channel.Attr("href")
		zhibos = append(zhibos, model.ZhiboProp{Name: name, URL: url})
	})
	return zhibos
}
