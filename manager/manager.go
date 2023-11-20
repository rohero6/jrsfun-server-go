package manager

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

var (
	Context playwright.BrowserContext
)

func InitManager() {
	pw, err := playwright.Run()
	if err != nil {
		log.Printf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Printf("could not launch browser: %v", err)
	}
	Context, err = browser.NewContext(playwright.BrowserNewContextOptions{
		// 设置为移动设备
		IsMobile: playwright.Bool(true),
		// 设置用户代理字符串（可选）
		UserAgent: playwright.String("Mozilla/5.0 (iPhone; CPU iPhone OS 13_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Mobile/15E148 Safari/604.1"),
	})
	if err != nil {
		log.Printf("could not create page: %v", err)
	}
	println("init playwright done")
}
