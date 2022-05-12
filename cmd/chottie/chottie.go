package main

import (
	"fmt"
	"github.com/xiaogogonuo/archimedes"
)

type ChineseHottie struct {
	archimedes.Archimedes
}

type BlogRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Query    string `json:"query"`
}

func (ch ChineseHottie) P1(response archimedes.Response) {
	text, err := response.Text()
	if err != nil {
		fmt.Println("P1", response.URL(), err.Error())
		return
	}
	fmt.Println(text)
}

func (ch ChineseHottie) Parser1(response archimedes.Response) {
	fmt.Println(response.Text())
}

func main() {
	ch := ChineseHottie{archimedes.New()}
	ch.Request("https://chottie.com/blog/categories", ch.P1, nil)
	//API := "https://chottie.com/api/mobile/search"
	//var br = BlogRequest{
	//	Page: 1,
	//	PageSize: 18,
	//	Query: "绯月樱",
	//}
	//body, _ := json.Marshal(&br)
	//request, _ := http.NewRequest(http.MethodPost, API, bytes.NewReader(body))
	//request.Header.Set("Content-Type", "application/json")
	//ch.NativeRequest(request, ch.P1, nil)
	ch.Boot()
}
