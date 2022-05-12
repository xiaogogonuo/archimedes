// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package archimedes

import (
	"net/http"
)

type archimedes struct {
	*engine
}

func (arch *archimedes) Request(url string, parser Parser, meta Meta) {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	arch.AdvanceRequest(&Request{Request: request, Parser: parser, Meta: meta, Filter: true})
}

func (arch *archimedes) NativeRequest(request *http.Request, parser Parser, meta Meta) {
	arch.AdvanceRequest(&Request{Request: request, Parser: parser, Meta: meta, Filter: true})
}

func (arch *archimedes) AdvanceRequest(request *Request) {
	go arch.engine.scheduler.advanceRequest(request)
}

func (arch *archimedes) Push(data interface{}) {

}

func (arch *archimedes) Boot() {
	arch.engine.start()
}

func New() Archimedes {
	return &archimedes{
		engine: newEngine(),
	}
}
