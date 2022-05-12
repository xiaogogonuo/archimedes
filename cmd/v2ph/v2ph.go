package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
)

// https://cdn.v2ph.com/photos/kbnXAJedAJLHGfE3.jpg

func TestNewSimulator() {
	sim := NewSimulator("chrome", "./chromedriver", 9999)
	opts := []selenium.ServiceOption{
		//selenium.Output(os.Stderr), // Output debug information to STDERR
	}
	service, err := sim.NewChromeService(opts...)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer service.Stop()

	// 禁止图片加载，加快渲染速度
	//imagCaps := map[string]interface{}{
	//	"profile.managed_default_content_settings.images": 2,
	//}

	chromeCaps := chrome.Capabilities{
		//Prefs: imagCaps,
		Path:  "",
		Args: []string{
			//"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			//"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}

	sim.AddChromeCap(chromeCaps)

	driver, err := sim.NewWebDriver()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer driver.Quit()
	defer driver.Close()

	if err := driver.Get("https://www.v2ph.com/actor/Fei-Yueying?hl=zh-Hant"); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 5)

	elements, err := driver.FindElements(selenium.ByXPATH, "//a[@class='media-cover']")
    var hrefs []string
	for _, element := range elements[:] {
		href, _ := element.GetAttribute("href")
		hrefs = append(hrefs, href)
		fmt.Println(href)
	}

	for _, href := range hrefs {
		_ = driver.Get(href)
		time.Sleep(time.Second * 5)
	}

	time.Sleep(time.Second * 10)
}

func main() {
	TestNewSimulator()
}
