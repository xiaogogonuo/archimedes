// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package archimedes

import (
	"bytes"
	"fmt"
	"github.com/xiaogogonuo/archimedes/pkg/encryption"
	"io"
	"net/http"
	"unsafe"
)

type archRequest struct {
	request *http.Request
	parser  Parser
	meta    Meta
	body    []byte // request body
	fb      []byte // request finger byte
	fs      string // request finger string
	client  *http.Client
}

func (req *archRequest) Valid() bool {
	return req.request != nil && req.request.URL != nil && req.parser != nil
}

func (req *archRequest) URL() string {
	return req.request.URL.String()
}

func (req *archRequest) Host() string {
	return req.request.Host
}

func (req *archRequest) Scheme() string {
	return req.request.URL.Scheme
}

func (req *archRequest) Method() string {
	return req.request.Method
}

func (req *archRequest) Body() []byte {
	if req.request.Body == nil {
		return nil
	}
	// TODO: see if there is an error exist or if there is other read method
	body, _ := io.ReadAll(req.request.Body)                // *http.Request will close Body after read
	req.request.Body = io.NopCloser(bytes.NewBuffer(body)) // rewrite body to *http.Request
	return body
}

func (req *archRequest) UMBByte() []byte {
	umb := make([]byte, 0)
	url := req.URL()
	method := req.Method()
	umb = append(umb, *(*[]byte)(unsafe.Pointer(&url))...)
	umb = append(umb, *(*[]byte)(unsafe.Pointer(&method))...)
	if req.body == nil {
		return umb
	}
	umb = append(umb, req.body...)
	return umb
}

func (req *archRequest) UMBString() string {
	umb := req.UMBByte()
	return *(*string)(unsafe.Pointer(&umb))
}

func (req *archRequest) FingerByte() []byte {
	return encryption.ArchHashByte(req.UMBByte(), "md5")
}

func (req *archRequest) FingerString() string {
	return fmt.Sprintf("%x", req.FingerByte())
}

type archResponse struct {
	response *http.Response
	*archRequest
	body      []byte // response body
	readCount uint
	readError error
}

func (res *archResponse) URL() string {
	return res.archRequest.URL()
}

func (res *archResponse) Meta() Meta {
	return res.archRequest.meta
}

func (res *archResponse) Byte() ([]byte, error) {
	if res.readCount == 0 {
		if res.response.Body == nil {
			res.readCount++
			res.body = nil
			res.readError = fmt.Errorf("%s: response body nill", res.URL())
			return nil, res.readError
		}
		// TODO: see if there is an error exist or if there is other read method
		body, err := io.ReadAll(res.response.Body)
		res.readCount++
		res.body = body
		res.readError = err
		_ = res.response.Body.Close()
		return body, err
	}
	return res.body, res.readError
}

func (res *archResponse) Text() (string, error) {
	body, err := res.Byte()
	return *(*string)(unsafe.Pointer(&body)), err
}
