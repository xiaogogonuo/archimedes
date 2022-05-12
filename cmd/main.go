// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/xiaogogonuo/archimedes"
	"net/http"
)

type App struct {
	archimedes.Archimedes
}

func (app App) Parser1(response archimedes.Response) {
	// 获取响应数据的字节流形式
	_, _ = response.Byte()
	// 获取响应数据的字符串形式
	_, _ = response.Text()
	// 获取请求携带过来的元数据
	_ = response.Meta()
	// 获取响应所对应的请求链接
	_ = response.URL()

	// 模拟发送从页面提取的新链接
	// Request方法是简易版的GET请求，客户端只需传入URL，自定义解析函数，元数据即可
	// 新链接：必填
	// 解析器：必填
	// 元数据：可选
	app.Request("https://www.douban.com", app.Parser2, map[string]interface{}{"parser": "Parser1"})
}

func (app App) Parser2(response archimedes.Response) {
	// 获取从Parser1传入的元数据
	_ = response.Meta()

	// 模拟发送从页面提取的新链接
	// NativeRequest方法是原生版的http请求，客户端需要传入自定义*http.Request、自定义解析函数、元数据
	// 新请求：必填
	// 解析器：必填
	// 元数据：可选
	customGetRequest, _ := http.NewRequest(http.MethodGet, "https://www.tencent.com", nil)
	app.NativeRequest(customGetRequest, app.Parser3, map[string]interface{}{"parser": "Parser2"})
}

func (app App) Parser3(response archimedes.Response) {
	// 获取从Parser2传入的元数据
	_ = response.Meta()

	// 模拟发送从页面提取的新链接
	// AdvanceRequest方法提供可选过滤器、自定义*http.Client
	// 新请求：必填
	// 解析器：必填
	// 元数据：可选
	// 过滤器：可选(应用场景：系统全局设置了过滤器，有部分特定的请求可以重复下载)
	// 客户端：可选
	req, _ := http.NewRequest(http.MethodGet, "https://www.tencent.com", nil)
	request := &archimedes.Request{
		Request: req,
		Parser:  app.Parser4,
		Meta:    map[string]interface{}{"parser": "Parser3"},
		Filter:  true, // false: 不过滤、true: 过滤 (该字段省略不写则代表不过滤，即默认不过滤)
		Client:  &http.Client{},
	}
	app.AdvanceRequest(request)
}

func (app App) Parser4(response archimedes.Response) {
	// 获取从Parser3传入的元数据
	_ = response.Meta()
}

func main() {
	app := App{archimedes.New()}
	app.Request("https://www.baidu.com", app.Parser1, nil)
	app.Boot()
}
