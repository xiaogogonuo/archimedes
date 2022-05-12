// Copyright 2022 JiaWei Lu <xiaogogonuo@163.com>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file.

package archimedes

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/pkg/errors"
	"github.com/xiaogogonuo/archimedes/pkg/counter"
	"sync"
)

const (
	DefaultFalsePositive     = 0.001
	DefaultEstimateItems     = 100000000
	DefaultRequestChanBuffer = 1 << 12
)

type scheduler struct {
	mu                sync.Mutex
	enableMapFilter   bool                // whether to use MapFilter
	enableBloomFilter bool                // whether to use BloomFilter
	fp                float64             // BloomFilter false positive
	ei                uint                // BloomFilter estimate items
	bf                *bloom.BloomFilter  // BloomFilter
	mf                map[string]struct{} // MapFilter
	allowedDomain     map[string]struct{} // allowed domain
	requestChan       chan *archRequest   // request channel
	ct                *counter.Counter
}

func (s *scheduler) advanceRequest(r *Request) {
	s.ct.IncrHandlingNumber()
	defer s.ct.DecrHandlingNumber()
	s.ct.IncrCalledCount()
	req := s.archRequestConstruct(r)
	if !s.archRequestValidityInspection(req) {
		s.ct.IncrInterceptCount()
		return
	}
	// 请求进行完有效性检验后，如果发现该请求不需要过滤，直接加入请求队列
	if !r.Filter {
		s.requestChan <- req
		s.ct.IncrAcceptedCount()
		return
	}
	if !s.archRequestRepeatabilityInspection(req, r.Filter) {
		s.ct.IncrInterceptCount()
		return
	}
	s.requestChan <- req
	s.ct.IncrAcceptedCount()
}

// archRequestConstruct 请求构造
func (s *scheduler) archRequestConstruct(r *Request) *archRequest {
	req := &archRequest{request: r.Request, parser: r.Parser, meta: r.Meta, client: r.Client}
	req.body = req.Body()
	req.fb = req.FingerByte()
	req.fs = req.FingerString()
	return req
}

// archRequestValidityInspection 请求有效性检验
func (s *scheduler) archRequestValidityInspection(req *archRequest) bool {
	if !req.Valid() {
		// TODO: use logger instead
		fmt.Println(req.URL(), "invalid req")
		return false
	}
	if !s.domainValid(req) {
		// TODO: use logger instead
		fmt.Println(req.URL(), "invalid domain")
		return false
	}
	if !s.schemeValid(req) {
		// TODO: use logger instead
		fmt.Println(req.URL(), "invalid scheme")
		return false
	}
	return true
}

// archRequestRepeatabilityInspection 请求重复性检验
func (s *scheduler) archRequestRepeatabilityInspection(req *archRequest, localFilter bool) bool {
	// 如果全局通用过滤器和全局布隆过滤器都不开启，则所有经过有效性检验的请求均可加入请求队列
	if !s.enableMapFilter && !s.enableBloomFilter {
		return true
	}
	// 如果全局布隆过滤器开启，不管全局通用过滤器是否开启，均使用全局布隆过滤器
	if s.enableBloomFilter {
		// 代表全局布隆过滤器里不存在此次请求
		if !s.bloomFilter(req) {
			return true
		}
		fmt.Printf("repeat request %s filter by bloom filter\n", req.URL())
		return false
	}
	// 代表全局通用过滤器里不存在此次请求
	if !s.mapFilter(req) {
		return true
	}
	fmt.Printf("repeat request %s filter by map filter\n", req.URL())
	return false
}

func (s *scheduler) domainValid(req *archRequest) bool {
	if len(s.allowedDomain) == 0 {
		return true
	}
	if _, ok := s.allowedDomain[req.Host()]; !ok {
		return false
	}
	return true
}

func (s *scheduler) schemeValid(req *archRequest) bool {
	if req.Scheme() == "http" || req.Scheme() == "https" {
		return true
	}
	return false
}

func (s *scheduler) mapFilter(req *archRequest) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.mf[req.fs]; !ok {
		s.mf[req.fs] = struct{}{}
		return false
	}
	return true
}

func (s *scheduler) bloomFilter(req *archRequest) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.bf.Test(req.fb) {
		s.bf.Add(req.fb)
		return false
	}
	return true
}

func (s *scheduler) SetAllowedDomains(domains ...string) {
	for _, domain := range domains {
		if _, ok := s.allowedDomain[domain]; !ok {
			s.allowedDomain[domain] = struct{}{}
		}
	}
}

func (s *scheduler) ResetMapFilter(enable bool) {
	if !enable {
		s.enableMapFilter = false
		s.mf = nil // release memory
		return
	}
	s.enableMapFilter = true
	s.mf = map[string]struct{}{}
}

func (s *scheduler) ResetBloomFilter(enable bool, ei uint, fp float64) {
	if !enable {
		s.enableBloomFilter = false
		s.bf = nil // release memory
		return
	}
	if fp > 1 || fp < 0 {
		// TODO: use logger instead
		panic(errors.New("bloom filter false positive should less than 1 and larger than 0"))
	}
	s.enableBloomFilter = true
	s.ei, s.fp = ei, fp
	s.bf = bloom.NewWithEstimates(s.ei, s.fp)
}

func (s *scheduler) SetRequestChanBuffer(n uint64) {
	s.requestChan = make(chan *archRequest, n)
}

func newScheduler() *scheduler {
	sd := &scheduler{
		enableMapFilter:   true,
		enableBloomFilter: true,
		ei:                DefaultEstimateItems,
		fp:                DefaultFalsePositive,
		mf:                map[string]struct{}{},
		allowedDomain:     map[string]struct{}{},
		requestChan:       make(chan *archRequest, DefaultRequestChanBuffer),
		ct:                &counter.Counter{},
	}
	sd.bf = bloom.NewWithEstimates(sd.ei, sd.fp)
	return sd
}
