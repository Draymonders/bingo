package group

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

// 不管是否有协程执行失败, wait()都要等待所有协程执行完成
func run(errF func() error) error {
	eg := errgroup.Group{}

	n := 10
	for i := 0; i < n; i++ {
		v := i
		eg.Go(func() error {
			return f(v, v, errF)
		})
	}
	return eg.Wait()
}

func runWithContext(ctx context.Context, errF func() error) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	n := 10
	for i := 0; i < n; i++ {
		v := i
		eg.Go(func() error {
			time.Sleep(3 * time.Millisecond)
			select {
			case <-ctx.Done():
				// ctx已经执行完了
				fmt.Printf("[%d] ctx.Done\n", v)
				return nil
			default:
				return f(v, v, errF)
			}
		})
	}

	return eg.Wait()
}

// id: 第多少个，t等多少秒
func f(id, t int, errF func() error) error {
	time.Sleep(time.Duration(t) * 10 * time.Millisecond)

	err := errF()
	if err != nil {
		fmt.Printf("[%d] err\n", id)
		return err
	}
	fmt.Printf("[%d] done\n", id)
	return errF()
}

func mockErr() error {
	return nil
}

func alwaysErr() error {
	return errors.New("233")
}

// 控制多少的失败，0 <= percentage <= 100
func randErr(percentage int) error {
	if percentage == 0 {
		return nil
	}
	n := 100
	if x := rand.Intn(n) + 1; x <= percentage {
		return errors.New("rand error")
	}
	return nil
}
