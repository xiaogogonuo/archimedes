package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

var html = `
<!DOCTYPE html>
<html>
<head></head>
<body>
    <div class="d1" alt="alt10">DIV10
        <a href="https://www.baidu.com">
        <div class="d2">
            <img src="https://123.jpg">
        </div>
    </div>
    <div class="d1" alt="alt11">DIV11
        <a href="https://www.tianmao.com">
        <div class="d2">
            <img src="https://456.jpg">
        </div>
    </div>
</body>
</html>
`
func main() {
	res, err := http.Get("https://www.baidu.com")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	b, err := httputil.DumpResponse(res, false)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(b))
	_ = res.Body.Close()

}
