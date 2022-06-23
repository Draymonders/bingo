package pool

import (
	"code.byted.org/gopkg/lang/maths"
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"
)

// todo @yubing 看看和 errGroup的区别

var token = struct{}{}

type BatchPool struct {
	batchLimit       chan struct{}
	allowPartSuccess bool
	eg               errgroup.Group
	err              atomic.Error
}

func NewBatchPool(batchNum int, allowPartSuccess bool) *BatchPool {
	batchNum = maths.MaxInt(1, batchNum)
	return &BatchPool{batchLimit: make(chan struct{}, batchNum), allowPartSuccess: allowPartSuccess}
}

func (b *BatchPool) Go(f func() error) {
	if b.isFinish() {
		return
	}

	b.batchLimit <- token

	b.eg.Go(func() (err error) {
		defer func() {
			<-b.batchLimit
		}()
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
