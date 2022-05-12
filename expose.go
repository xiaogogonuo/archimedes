// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package archimedes

import (
	"net/http"
	"time"
)

// Request is a scalable request
type Request struct {
	*http.Request
	Parser
	Meta
	Filter bool
	*http.Client
}

type Archimedes interface {
	// Request is a widely used http GET request
	Request(url string, parse Parser, meta Meta)

	// NativeRequest is a user defined http request with optional header、body
	NativeRequest(request *http.Request, parse Parser, meta Meta)

	// AdvanceRequest provide more advanced function
	AdvanceRequest(request *Request)

	// Push send data to database
	Push(interface{})

	// Boot start the server
	Boot()

	// Setter is an optional setting for client
	Setter

	// Logger is an optional for client
	//Logger
}

type Response interface {
	URL() string
	Byte() ([]byte, error)
	Text() (string, error)
	Meta() Meta
}

type Parser func(Response)

type Meta map[string]interface{}

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
	Panic(string)
	Fatal(string)
	DebugF(string, ...interface{})
	InfoF(string, ...interface{})
	WarnF(string, ...interface{})
	ErrorF(string, ...interface{})
	PanicF(string, ...interface{})
	FatalF(string, ...interface{})
}

type Setter interface {
	// SetRequestChanBuffer set the buffer of request channel
	SetRequestChanBuffer(uint64)

	// SetDownloadConcurrent set the download concurrent
	SetDownloadConcurrent(uint64)

	// SetDownloadDelay set the download delay
	SetDownloadDelay(time.Duration)

	// SetDownloadTimeout set the download timeout
	SetDownloadTimeout(time.Duration)

	// SetRetryEnable set whether to retry after download failure, default enabled
	SetRetryEnable(bool)

	// SetRetryTimes set the number of retries
	SetRetryTimes(uint64)

	// SetRetryHttpCodes set the status code that needs to be downloaded again
	SetRetryHttpCodes(...int)

	// SetAllowedDomains set allowed domain, default allows all domains
	SetAllowedDomains(...string)

	// ResetMapFilter reset map filter
	ResetMapFilter(bool)

	// ResetBloomFilter reset bloom filter with new estimated items and false positive
	// SetBloomFilterEnable set whether to use bloom filter, default enabled
	// if BloomFilter is enabled, whether MapFilter enabled or not, BloomFilter will be used
	// if BloomFilter is disabled, MapFilter is enabled, MapFilter will be used
	// if BloomFilter and MapFilter all disabled, no Filter will be used
	// if estimated urls is larger than 1e8, BloomFilter recommended
	ResetBloomFilter(bool, uint, float64)
}
