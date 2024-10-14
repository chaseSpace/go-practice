package main

import (
	"encoding/json"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"os"
	"time"
)

func main() {
	buf, _ := os.ReadFile("cookies.json")

	// Launch a new browser with default options, and connect to it.
	//browser := rod.New().MustConnect()
	//
	//// Even you forget to close, rod will close it after main process ends.
	//defer browser.MustClose()
	//
	//// Create a new page
	//page := browser.MustPage("https://www.instagram.com/misterparadisenyc").MustWaitStable()
	//s, _ := page.HTML()
	//
	//println(s)

	l := launcher.New().
		Headless(false).
		Devtools(true)

	defer l.Cleanup()

	url := l.MustLaunch()

	// Trace shows verbose debug information for each action executed
	// SlowMotion is a debug related function that waits 2 seconds between
	// each action, making it easier to inspect what your code is doing.
	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(time.Millisecond * 200).
		MustConnect()

	// ServeMonitor plays screenshots of each tab. This feature is extremely
	// useful when debugging with headless mode.
	// You can also enable it with flag "-rod=monitor"
	launcher.Open(browser.ServeMonitor(""))

	defer browser.MustClose()

	//https://www.instagram.com/misterparadisenyc
	var ck []*proto.NetworkCookie
	_ = json.Unmarshal(buf, &ck)
	browser.MustSetCookies(ck...)

	page := browser.MustPage("https://www.instagram.com/misterparadisenyc").
		MustWaitDOMStable()
	//page.MustElement("#loginForm > div > div:nth-child(1) > div > label > input").
	//	MustInput("leigege_ya")
	//page.MustElement("#loginForm > div > div:nth-child(2) > div > label > input").
	//	MustInput("1213is..").MustType(input.Enter)
	//
	//page.MustElementX("//span[text()='Home']").MustWaitStable()
	//
	//ck, _ = page.Cookies(nil)
	//println(len(ck), "cookies")
	//buf2, _ := json.Marshal(ck)
	//_ = os.WriteFile("cookies.json", buf2, 0666)

	//page.MustSetCookies()

	//var err error
	//elem := page.MustElementX("//section/main/div/div[2]/div[1]/div/div")

	baseTags := page.MustElementsX("//section/main/div/div[2]/div/div//a[@href]")
	imgs := page.MustElementsX("//section/main/div/div[2]/div/div//a[@href]//img")

	if len(baseTags) != len(imgs) {
		println("len(baseTags) != len(imgs)", len(baseTags), len(imgs))
		return
	}

	for i, e := range baseTags {
		//println(11111, i, e.MustHTML())
		//fmt.Println("序号", i+1, "内容", *e.MustAttribute("alt"))
		//fmt.Println("序号", i+1, "图片", *e.MustAttribute("src"))
		println(1111, i+1, *e.MustAttribute("href"),
			*imgs[i].MustAttribute("alt"),
			"xxx",
			*imgs[i].MustAttribute("src"),
		)
		//e.MustClick()
		//println("333...", e.MustElementX("//ul/li//div[@class='_aagv']//img").MustAttribute("src"))
		//e.MustType(input.Escape)
	}
	println(66666)
	time.Sleep(time.Hour)
}
