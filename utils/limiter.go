package utils

import (
	"sync/atomic"
	"time"
)

var _ Limiter = &QpsLimiter{}

type Limiter interface {
	// Acquire 获取令牌
	Acquire() bool

	Stop()
}

type QpsLimiter struct {
	tokens int32
	limit  int32
	stopCh chan struct{}

	// interval 秒可以发limit个令牌，由于单机的，limit是int32类型就可以
	interval time.Duration
}

func NewQpsLimiter(interval time.Duration, limit int32) Limiter {
	l := &QpsLimiter{
		tokens:   limit,
		limit:    limit,
		interval: interval,
		stopCh:   make(chan struct{}, 1),
	}
	go l.runTicker()
	return l
}

// 1秒发多少个tokens
func calc(n time.Duration, m int) int {
	if n < time.Second {
		n = time.Second
	}
	return m / (int(n) / int(time.Second))
}

func (q *QpsLimiter) Acquire() bool {
	tokens := atomic.LoadInt32(&q.tokens)
	for tokens > 0 {
		if atomic.CompareAndSwapInt32(&q.tokens, tokens, tokens-1) {
			return true
		}
		tokens = atomic.LoadInt32(&q.tokens)
	}
	if tokens <= 0 {
		return false
	}
	return true
}

func (q *QpsLimiter) Stop() {
	q.stopCh <- struct{}{}
}

func (q *QpsLimiter) runTicker() {
	ticker := time.NewTicker(q.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			q.updateToken()
		case <-q.stopCh:
			return
		}
	}
}

func (q *QpsLimiter) updateToken() {
	tokens, limit := atomic.LoadInt32(&q.tokens), atomic.LoadInt32(&q.limit)

	if tokens > limit {
		return
	}
	atomic.StoreInt32(&q.tokens, limit)
}
