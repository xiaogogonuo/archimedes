// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package archimedes

import (
	"context"
	"fmt"
	"time"
)

const (
	DefaultIdleCount = 10
	DefaultHeartbeat = time.Second
)

var ctx, cancel = context.WithCancel(context.Background())

type engine struct {
	idleCount int
	heartbeat time.Duration
	done      chan int
	*scheduler
	*downloader
}

func (e *engine) dispatch() {
	go func() {
		for {
			select {
			case req := <-e.scheduler.requestChan:
				e.downloader.download(req, e.scheduler.requestChan)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (e *engine) idle() bool {
	if e.scheduler.ct.HandlingNumber() == 0 &&
		e.downloader.ct.HandlingNumber() == 0 {
		return true
	}
	return false
}

func (e *engine) monitor() {
	go func() {
		var count int
		defer func() {
			cancel()
			e.done <- count
		}()
		for range time.Tick(e.heartbeat) {
			if e.idle() {
				count++
			}
			if count > e.idleCount {
				if e.idle() {
					break
				} else {
					count = 0
				}
			}
		}
	}()
}

func (e *engine) summary() {
	fmt.Println()
	fmt.Println("* * * * * * * * * * * * * * * * 统计信息 * * * * * * * * * * * * * * * *")
	fmt.Println("客户端发起请求的总数量 = 客户端请求被拦截的数量 + 客户端请求被接受的数量")
	fmt.Println("客户端请求被接受的数量 = 客户端请求下载失败数量 + 客户端请求下载成功数量")
	fmt.Println()
	fmt.Printf("客户端发起请求的总数量: %d个\n", e.scheduler.ct.CalledCount())
	fmt.Printf("客户端请求被拦截的数量: %d个\n", e.scheduler.ct.InterceptCount())
	fmt.Printf("客户端请求被接受的数量: %d个\n", e.scheduler.ct.AcceptedCount())
	fmt.Printf("客户端请求下载失败数量: %d个\n", e.downloader.ct.FailedCount())
	fmt.Printf("客户端请求下载成功数量: %d个\n", e.downloader.ct.CompletedCount())
	fmt.Println("* * * * * * * * * * * * * * * * 统计信息 * * * * * * * * * * * * * * * *")
}

func (e *engine) start() {
	e.dispatch()
	e.monitor()
	<-e.done
	e.summary()
}

func newEngine() *engine {
	return &engine{
		idleCount:  DefaultIdleCount,
		heartbeat:  DefaultHeartbeat,
		done:       make(chan int),
		scheduler:  newScheduler(),
		downloader: newDownloader(),
	}
}
