package usage2http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// URI通用结构描述 scheme:[//[user:password@]host[:port]][/]path[?query][#fragment]

func ParseRequestURI() {
	api := "https://127.0.0.1:8080/rest/report?id=1&name=work#we"
	reqURL, _ := url.ParseRequestURI(api) // *url.URL, error

	fmt.Println(reqURL.Scheme) // https
	fmt.Println(reqURL.Opaque)
	fmt.Println(reqURL.User)
	fmt.Println(reqURL.Host) // 127.0.0.1:8080
	fmt.Println(reqURL.Path) // /rest/report
	fmt.Println(reqURL.RawPath)
	fmt.Println(reqURL.ForceQuery) // false
	fmt.Println(reqURL.RawQuery) // id=1&name=work
	fmt.Println(reqURL.Fragment) // we
	fmt.Println(reqURL.RawFragment)

	// 获取请求参数
	values := reqURL.Query() // map[id:[1] name:[work]]
	// 添加请求参数
	values.Add("month", "20220430")
	reqURL.RawQuery = values.Encode()
	fmt.Println(reqURL.String()) // https://127.0.0.1:8080/rest/report?id=1&month=20220430&name=work%23we
}

func NewRequestURI() {
	api := "https://127.0.0.1:8080/rest/report?id=1&name=work#we"
	req, _ := http.NewRequest(http.MethodGet, api, nil)
	reqURL := req.URL // *url.URL

	fmt.Println(reqURL.Scheme) // https
	fmt.Println(reqURL.Opaque)
	fmt.Println(reqURL.User)
	fmt.Println(reqURL.Host) // 127.0.0.1:8080
	fmt.Println(reqURL.Path) // /rest/report
	fmt.Println(reqURL.RawPath)
	fmt.Println(reqURL.ForceQuery) // false
	fmt.Println(reqURL.RawQuery) // id=1&name=work
	fmt.Println(reqURL.Fragment) // we
	fmt.Println(reqURL.RawFragment)

	// 获取请求参数
	values := reqURL.Query() // map[id:[1] name:[work]]
	// 添加请求参数
	values.Add("month", "20220430")
	reqURL.RawQuery = values.Encode()
	fmt.Println(reqURL.String()) // https://127.0.0.1:8080/rest/report?id=1&month=20220430&name=work#we
}

func Get() {
	// The timeout includes connection time, any redirects, and reading the response body
	client := http.Client{Timeout: time.Millisecond} //
	req, _ := http.NewRequest(http.MethodGet, "https://www.baidu.com", nil)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(io.ReadAll(res.Body))
	_ = res.Body.Close()
}

func ResponseBody() {
	res, err := http.Get("https://www.baidu.com")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(b))

	_ = res.Body.Close()
}

func ResponseAll() {
	res, err := http.Get("https://www.baidu.com")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	b, err := httputil.DumpResponse(res, true) // set false to exclude body
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(b))

	_ = res.Body.Close()

	/*
	HTTP/1.1 200 OK
	Content-Length: 227
	Accept-Ranges: bytes
	Cache-Control: no-cache
	Connection: keep-alive
	Content-Type: text/html
	Date: Sat, 30 Apr 2022 04:53:07 GMT
	P3p: CP=" OTI DSP COR IVA OUR IND COM "
	P3p: CP=" OTI DSP COR IVA OUR IND COM "
	Pragma: no-cache
	Server: BWS/1.1
	Set-Cookie: BD_NOT_HTTPS=1; path=/; Max-Age=300
	Set-Cookie: BIDUPSID=5A9CC2B34C43C02449076C9E204E051E; expires=Thu, 31-Dec-37 23:55:55 GMT; max-age=2147483647; path=/; domain=.baidu.com
	Set-Cookie: PSTM=1651294387; expires=Thu, 31-Dec-37 23:55:55 GMT; max-age=2147483647; path=/; domain=.baidu.com
	Set-Cookie: BAIDUID=5A9CC2B34C43C024CC3C21AAD1FCF990:FG=1; max-age=31536000; expires=Sun, 30-Apr-23 04:53:07 GMT; domain=.baidu.com; path=/; version=1; comment=bd
	Strict-Transport-Security: max-age=0
	Traceid: 1651294387026665319410622539328100134642
	X-Frame-Options: sameorigin
	X-Ua-Compatible: IE=Edge,chrome=1

	<html>
	<head>
	        <script>
	                location.replace(location.href.replace("https://","http://"));
	        </script>
	</head>
	<body>
	        <noscript><meta http-equiv="refresh" content="0;url=http://www.baidu.com/"></noscript>
	</body>
	</html>
	*/
}

//func GetWithoutBody() {
//	request, _ := http.NewRequest(http.MethodGet, url, nil)
//	request.Header.Set("Cookie", "a=1; b=2")
//	client := http.Client{}
//	response, _ := client.Do(request)
//	fmt.Println(response)
//}
