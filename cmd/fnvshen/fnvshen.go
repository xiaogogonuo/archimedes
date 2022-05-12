package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/xiaogogonuo/archimedes"
	"io"
	"os"
	"path"
	"strings"
)

const (
	domain = "https://m.fnvshen.com"
)

type FNvShen struct {
	archimedes.Archimedes
}

func (f FNvShen) Parser1(response archimedes.Response) {
	text, _ := response.Text()
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(text))
	dom.Find("#gallerydiv a").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok {
			f.Request(href, f.Parser2, nil)
		}
	})
	dom.Find("a[class='next']").Each(func(i int, selection *goquery.Selection) {
		if next, ok := selection.Attr("href"); ok {
			f.Request(next, f.Parser1, nil)
		}
	})
}

func (f FNvShen) Parser2(response archimedes.Response) {
	text, _ := response.Text()
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(text))
	// 获取用户名称
	var name string
	dom.Find(".ck-title-set span").Each(func(i int, selection *goquery.Selection) {
		name = selection.Text()
	})
	// 获取用户首页
	dom.Find("#dgirl a").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok {
			href = domain + href + "album/"
			f.Request(href, f.Parser4, nil)
		}
	})
	// 获取下一页
	dom.Find(".next").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		href = domain + href
		f.Request(href, f.Parser2, nil)
	})
	// 获取图片
	dom.Find("#idiv img").Each(func(i int, selection *goquery.Selection) {
		src, _ := selection.Attr("src")
		alt, _ := selection.Attr("alt")
		if strings.Contains(src, "t1.buuxk.com") {
			src = strings.ReplaceAll(src, "t1.buuxk.com", "img.buuxk.com")
		}
		f.Request(src, f.Parser3, map[string]interface{}{"alt": alt, "name": name})
	})
}

func (f FNvShen) Parser3(response archimedes.Response) {
	ul := strings.Split(response.URL(), "/")
	suf := ul[len(ul)-1]
	alt := response.Meta()["alt"].(string)
	alts := strings.Split(alt, " ")
	alts = alts[:len(alts)-1]
	alt = strings.Join(alts, "_")
	savePath := path.Join("/Users/lujiawei/Documents/golang/go-wallpaper", "image",
		response.Meta()["name"].(string), alt)
	_ = CreateMultiDir(savePath)
	body, _ := response.Byte()
	if len(body) < 10000 {
		fmt.Println(response.URL(), "image too small, re download")
		f.Request(response.URL(), f.Parser3, response.Meta())
		return
	}
	file, _ := os.OpenFile(path.Join(savePath, suf), os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0777)
	_, _ = io.Copy(file, bytes.NewReader(body))
}

func (f FNvShen) Parser4(response archimedes.Response) {
	text, _ := response.Text()
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(text))

	dom.Find("#dphoto a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		href = domain + href
		f.Request(href, f.Parser2, nil)
	})

	// 获取下一页
	dom.Find(".next").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		href = domain + href
		f.Request(href, f.Parser4, nil)
	})
}

func main() {
	ns := FNvShen{archimedes.New()}
	//ns.SetFilterEnable(false)
	//ns.SetBloomFilterEnable(false)
	ns.Request("https://m.fnvshen.com/gallery/", ns.Parser1, nil)
	ns.Boot()
}
