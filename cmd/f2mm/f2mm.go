package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/xiaogogonuo/archimedes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	domain = "https://www.f2mm.com"
	base = "/Users/lujiawei/Documents/golang/go-wallpaper/image"
)

type F2MM struct {
	archimedes.Archimedes
}

func (f *F2MM) Parser1(response archimedes.Response) {
	text, _ := response.Text()
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(text))
	dom.Find("div[class='card-body row'] a").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok {
			f.Request(href, f.Parser2, nil)
		}
	})
	dom.Find("a[aria-label='Next »']").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok {
			f.Request(href, f.Parser1, nil)
		}
	})
}

func (f *F2MM) Parser2(response archimedes.Response) {
	text, _ := response.Text()
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(text))
	var title string
	dom.Find("h1[class='post-title']").Each(func(i int, selection *goquery.Selection) {
		title = selection.Text()
	})
	dom.Find("#masonry div").Each(func(i int, selection *goquery.Selection) {
		if src, ok := selection.Attr("data-src"); ok {
			if strings.HasSuffix(src, "holder.png") {
				return
			}
			f.Request(src, f.Parser3, map[string]interface{}{"title": title})
		}
	})
}

func (f *F2MM) Parser3(response archimedes.Response) {
	fmt.Println(response.URL(), response.Meta())
	binary, _ := response.Byte()
	if len(binary) < 10000 {
		fmt.Println(response.URL(), "image too small, re download")
		request, _ := http.NewRequest(http.MethodGet, response.URL(), nil)
		f.AdvanceRequest(&archimedes.Request{
			Request: request,
			Parser: f.Parser3,
			Meta: response.Meta(),
			Filter: false,
		})
		return
	}
	url := strings.Split(response.URL(), "/")
	suf := url[len(url)-1]
	title := response.Meta()["title"].(string)
	ts := strings.Split(title, " ")
	var name string
	var save string
	if len(ts) < 4 {
		save = path.Join(base, title)
	} else {
		name = ts[3]
		save = path.Join(base, name, strings.Join(ts[:3], " "))
	}
	_ = CreateMultiDir(save)
	file, _ := os.OpenFile(path.Join(save, suf), os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0777)
	_, _ = io.Copy(file, bytes.NewReader(binary))
}

func main() {
	f := &F2MM{archimedes.New()}
	f.Request("https://www.f2mm.com/beauty", f.Parser1, nil)
	f.Boot()
}
