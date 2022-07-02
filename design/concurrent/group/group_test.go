package group

import (
	"context"
	"fmt"
	"testing"
)

func Test_WaitOne(t *testing.T) {
	waitOne()
}

func Test_ErrGroup(t *testing.T) {
	{
		if err := run(mockErr); err != nil {
			t.Logf("mockErr!!")
			t.FailNow()
		}
	}

	{
		fmt.Println("=====")
		if err := run(alwaysErr); err == nil {
			t.Logf("alwaysErr, but not error!!")
			t.FailNow()
		} else {
			t.Logf("always: %v", err)
		}
	}
	ctx := context.Background()
	{
		fmt.Println("=====")
		runWithContext(ctx, alwaysErr)
	}
}

func Test_BatchPool(t *testing.T) {
	{
		p := NewBatchPool(2, false)
		n := 10
		for i := 0; i < n; i++ {
			v := i
			p.Go(func() error {
				return f(v, v, func() error {
					return alwaysErr()
				})
			})
		}

		err := p.Wait()
		if err == nil {
			t.FailNow()
		}
	}
}
