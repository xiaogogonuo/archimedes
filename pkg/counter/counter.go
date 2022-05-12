// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package counter

import (
	"sync/atomic"
)

type Counter struct {
	calledCount    uint64 // 代表请求调用的计数
	failedCount    uint64 // 代表请求失败的计数
	acceptedCount  uint64 // 代表请求被接受的计数
	interceptCount uint64 // 代表请求被拦截的计数
	completedCount uint64 // 代表请求成功完成的计数
	handlingNumber uint64 // 代表请求实时处理的计数
}

func (m *Counter) IncrCalledCount() {
	atomic.AddUint64(&m.calledCount, 1)
}

func (m *Counter) DecrCalledCount() {
	atomic.AddUint64(&m.calledCount, ^uint64(0))
}

func (m *Counter) IncrFailedCount() {
	atomic.AddUint64(&m.failedCount, 1)
}

func (m *Counter) DecrFailedCount() {
	atomic.AddUint64(&m.failedCount, ^uint64(0))
}

func (m *Counter) IncrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, 1)
}

func (m *Counter) DecrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, ^uint64(0))
}

func (m *Counter) IncrInterceptCount() {
	atomic.AddUint64(&m.interceptCount, 1)
}

func (m *Counter) DecrInterceptCount() {
	atomic.AddUint64(&m.interceptCount, ^uint64(0))
}

func (m *Counter) IncrCompletedCount() {
	atomic.AddUint64(&m.completedCount, 1)
}

func (m *Counter) DecrCompletedCount() {
	atomic.AddUint64(&m.completedCount, ^uint64(0))
}

func (m *Counter) IncrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, 1)
}

func (m *Counter) DecrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, ^uint64(0))
}

func (m *Counter) CalledCount() uint64 {
	return atomic.LoadUint64(&m.calledCount)
}

func (m *Counter) FailedCount() uint64 {
	return atomic.LoadUint64(&m.failedCount)
}

func (m *Counter) AcceptedCount() uint64 {
	return atomic.LoadUint64(&m.acceptedCount)
}

func (m *Counter) InterceptCount() uint64 {
	return atomic.LoadUint64(&m.interceptCount)
}

func (m *Counter) CompletedCount() uint64 {
	return atomic.LoadUint64(&m.completedCount)
}

func (m *Counter) HandlingNumber() uint64 {
	return atomic.LoadUint64(&m.handlingNumber)
}

func (m *Counter) Clear() {
	atomic.StoreUint64(&m.calledCount, 0)
	atomic.StoreUint64(&m.failedCount, 0)
	atomic.StoreUint64(&m.acceptedCount, 0)
	atomic.StoreUint64(&m.interceptCount, 0)
	atomic.StoreUint64(&m.completedCount, 0)
	atomic.StoreUint64(&m.handlingNumber, 0)
}
