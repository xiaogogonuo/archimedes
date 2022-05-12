// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package archimedes

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/xiaogogonuo/archimedes/pkg/counter"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultRetryTimes         = 10
	DefaultDownloadDelay      = time.Microsecond
	DefaultDownloadTimeout    = time.Second * 30
	DefaultDownloadConcurrent = 1 << 10
)

var DefaultRetryHttpCodes = map[int]struct{}{
	http.StatusRequestTimeout:      {},
	http.StatusTooManyRequests:     {},
	http.StatusInternalServerError: {},
	http.StatusBadGateway:          {},
	http.StatusServiceUnavailable:  {},
	http.StatusGatewayTimeout:      {},
}

type downloader struct {
	mu             sync.Mutex
	delay          time.Duration
	timeout        time.Duration
	retryTimes     uint64
	retryEnable    bool
	retryHttpCodes map[int]struct{}
	retryCounter   map[string]uint64
	concurrent     chan struct{}
	ct             *counter.Counter
}

func (d *downloader) download(req *archRequest, reqChan chan *archRequest) {
	d.concurrent <- struct{}{}
	go d.downloading(req, reqChan)
}

func (d *downloader) downloading(req *archRequest, reqChan chan *archRequest) {
	d.ct.IncrHandlingNumber()
	defer d.ct.DecrHandlingNumber()
	defer func() {
		time.Sleep(d.delay)
		<-d.concurrent
	}()
	response := &archResponse{}
	done := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	go func() {
		client := d.archClient()
		if req.client != nil {
			client = req.client
		}
		res, err := client.Do(req.request)
		response.response = res
		response.archRequest = req
		done <- err
	}()
	select {
	case err := <-done:
		//if err != nil {
		//	// TODO: use logger instead
		//	fmt.Println(err.Error(), req.URL())
		//	d.ct.IncrFailedCount()
		//	return
		//}


		if err != nil {
			if !d.retryEnable {
				d.ct.IncrFailedCount()
				return
			}
			if !d.ifRetry(req) {
				// TODO: use logger instead
				fmt.Println(req.URL(), "up to max retry times, drop it")
				d.ct.IncrFailedCount()
				return
			}
			go func() { reqChan <- req }()
			return
		}

		// retry enable points at retry http code
		if d.isRetryHttpCode(response.response) {
			if !d.retryEnable {
				d.ct.IncrFailedCount()
				return
			}
			if !d.ifRetry(req) {
				// TODO: use logger instead
				fmt.Println(req.URL(), "up to max retry times, drop it")
				d.ct.IncrFailedCount()
				return
			}
			go func() { reqChan <- req }()
			return
		}
		d.ct.IncrCompletedCount()
		go req.parser(response)
	case <-ctx.Done():
		if !d.retryEnable {
			d.ct.IncrFailedCount()
			return
		}
		// download timeout process mechanism same as retry http code
		if !d.ifRetry(req) {
			// TODO: use logger instead
			fmt.Println(req.URL(), "up to max retry times, drop it")
			d.ct.IncrFailedCount()
			return
		}
		go func() { reqChan <- req }()
	}
}

func (d *downloader) ifRetry(req *archRequest) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	count, ok := d.retryCounter[req.fs]
	if !ok {
		d.retryCounter[req.fs] = 1
		return true
	}
	if count > d.retryTimes {
		return false
	}
	d.retryCounter[req.fs]++
	return true
}

func (d *downloader) isRetryHttpCode(response *http.Response) bool {
	if _, ok := d.retryHttpCodes[response.StatusCode]; ok {
		return true
	}
	return false
}

func (d *downloader) archClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (d *downloader) SetDownloadDelay(delay time.Duration) {
	d.delay = delay
}

func (d *downloader) SetDownloadTimeout(timeout time.Duration) {
	d.timeout = timeout
}

func (d *downloader) SetDownloadConcurrent(concurrent uint64) {
	d.concurrent = make(chan struct{}, concurrent)
}

func (d *downloader) SetRetryEnable(enable bool) {
	d.retryEnable = enable
}

func (d *downloader) SetRetryTimes(times uint64) {
	d.retryTimes = times
}

func (d *downloader) SetRetryHttpCodes(code ...int) {
	for _, c := range code {
		if _, ok := d.retryHttpCodes[c]; !ok {
			d.retryHttpCodes[c] = struct{}{}
		}
	}
}

func newDownloader() *downloader {
	return &downloader{
		delay:          DefaultDownloadDelay,
		timeout:        DefaultDownloadTimeout,
		retryTimes:     DefaultRetryTimes,
		retryEnable:    true,
		retryHttpCodes: DefaultRetryHttpCodes,
		retryCounter:   map[string]uint64{},
		concurrent:     make(chan struct{}, DefaultDownloadConcurrent),
		ct:             &counter.Counter{},
	}
}
