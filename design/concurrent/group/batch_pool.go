package group

import (
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"
)

type BatchPool struct {
	allowPartSuccess bool
	eg               *errgroup.Group
	err              atomic.Error
}

// allowPartSuccess: true, 所有都执行
// allowPartSuccess: false, 遇到失败，后续goroutine 快速关闭
func NewBatchPool(batchNum int, allowPartSuccess bool) *BatchPool {
	var eg errgroup.Group
	if batchNum > 1 {
		eg.SetLimit(batchNum)
	}
	return &BatchPool{eg: &eg, allowPartSuccess: allowPartSuccess}
}

func (b *BatchPool) Go(f func() error) {
	if b.isFinish() {
		return
	}

	b.eg.Go(func() (err error) {
		if b.isFinish() { // fast fail 只要有一个失败就不继续走下去了
			return
		}
		if err = f(); err != nil {
			b.err.Store(err)
		}
		return
	})
}

func (b *BatchPool) Wait() error {
	_ = b.eg.Wait()
	err := b.err.Load()
	b.err.Store(error(nil))

	return err
}

func (b *BatchPool) isFinish() bool {
	return b.err.Load() != nil && !b.allowPartSuccess
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
