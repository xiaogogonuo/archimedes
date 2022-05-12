package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/xiaogogonuo/archimedes"
	"github.com/xiaogogonuo/archimedes/pkg/logger"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// GroupNews write to table `t_dmbe_group_news_info`
type GroupNews struct {
	NewsGuid       string   `json:"newsGuid"`       // 新闻主键
	NewsTitle      string   `json:"newsTitle"`      // 新闻标题
	NewsTs         string   `json:"newsTs"`         // 新闻日期
	NewsUrl        string   `json:"newsUrl"`        // 新闻链接
	NewsSource     string   `json:"newsSource"`     // 新华网
	NewsSourceCode string   `json:"newsSourceCode"` // WEB_00053
	NewsSummary    string   `json:"newsSummary"`    // 新闻正文
	PolicyType     string   `json:"policyType"`     // 10
	PolicyTypeName string   `json:"policyTypeName"` // 国家政策
	NewsGysCode    string   `json:"newsGysCode"`    // 90
	NewsGysName    string   `json:"newsGysName"`    // 爬虫
	NewsId         int      `json:"newsId"`         // 0
	Image          [][]byte `json:"image"`          // 新闻图片
}

const WebService = "http://106.37.165.121/inf/chengtong/py/sy/groupNewsInfo/saveGroupNewsInfo"

func Send(api string, data GroupNews) {
	postData := map[string][]GroupNews{"data": {data}}
	m, _ := json.Marshal(postData)
	req, err := http.NewRequest(http.MethodPost, api, bytes.NewReader(m))
	if err != nil {
		logger.Error(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	logger.Info(string(b))
}

type App struct {
	mu sync.Mutex
	m  map[string]GroupNews
	archimedes.Archimedes
}

func (app App) Parser1(response archimedes.Response) {
	text, _ := response.Text()
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(text))
	// 财经24小时Top5
	dom.Find("div[class='cjtop'] a").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok {
			app.Request(href, app.Parser2, nil)
		}
	})

	// 滚动
	// dom.Find("ul[class='silder_nav clearfix'] a").Each(func(i int, selection *goquery.Selection) {
	// 	if href, ok := selection.Attr("href"); ok {
	// 		app.Request(href, app.Parser2, nil)
	// 	}
	// })
	// dom.Find("div[class='tit']>a").Each(func(i int, selection *goquery.Selection) {
	// 	if href, ok := selection.Attr("href"); ok {
	// 		app.Request(href, app.Parser2, nil)
	// 	}
	// })
}

func (app App) Parser2(response archimedes.Response) {
	url := response.URL()
	text, _ := response.Text()
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(text))
	var title string
	dom.Find("span[class='title']").Each(func(i int, selection *goquery.Selection) {
		text := strings.Trim(selection.Text(), "\n\t ")
		title = text
	})
	if title == "" {
		return
	}
	// 新闻正文
	var texts []string
	dom.Find("#detail>p").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		texts = append(texts, text)
	})

	var publishDate string
	dom.Find("span[class='time']").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		publishDate = strings.Trim(text, " \n\t")
	})	

	date := strings.Join(strings.Split(url, "/")[4:6], "-")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	date += (" " + publishDate)
	var gn GroupNews
	hash := md5.New()
	hash.Write(*(*[]byte)(unsafe.Pointer(&url)))
	gn.NewsGuid = fmt.Sprintf("%x", hash.Sum(nil))
	gn.NewsTitle = title
	gn.NewsTs = date
	gn.NewsUrl = response.URL()
	gn.NewsSource = "新华网"
	gn.NewsSourceCode = "WEB_00053"
	gn.NewsSummary = strings.Join(texts, "")
	gn.PolicyType = "10"
	gn.PolicyTypeName = "国家政策"
	gn.NewsGysCode = "90"
	gn.NewsGysName = "爬虫"
	gn.NewsId = 0
	app.mu.Lock()
	if _, ok := app.m[gn.NewsGuid]; !ok {
		app.m[gn.NewsGuid] = gn
	}
	app.mu.Unlock()
}

func main() {
	app := App{
		m:          map[string]GroupNews{},
		Archimedes: archimedes.New(),
	}
	app.SetAllowedDomains("www.news.cn", "www.xinhuanet.com")
	app.Request("http://www.xinhuanet.com/fortunepro/", app.Parser1, nil)
	app.Boot()
	for _, gn := range app.m {
		// fmt.Println(gn.NewsTs)
		Send(WebService, gn)
	}
}

// GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o xinhua cmd/xinhua/xinhua.go
