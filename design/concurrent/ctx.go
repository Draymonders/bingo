package concurrent

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/draymonders/bingo/utils"
)

const (
	ctxCaseKey = "case_key"
)

var errTimeout = errors.New("timeout")

type Job func(ctx context.Context)

// 实现N个任务并发跑，并且用 ctx 控制超时时间，超时的话，直接返回
// 2022.06.21
func ParallelRun(ctx context.Context, jobs []Job) error {
	if len(jobs) == 0 {
		return nil
	}
	nums := len(jobs)
	wg := sync.WaitGroup{}
	fmt.Printf("[ctx %v] len(jobs)=%v\n", ctx.Value(ctxCaseKey), nums)

	done := make(chan struct{}, 1)
	// 不用加defer，加了的话，会引起panic，由gc来释放channel
	// defer close(done)

	wg.Add(nums)
	for i := 0; i < nums; i++ {
		go func(job Job) {
			defer wg.Done()
			job(ctx)
		}(jobs[i])
	}

	go func() {
		defer utils.Recover("all jobs done")
		wg.Wait()

		select {
		case done <- struct{}{}:
		}

	}()

	// 用time.After也可实现
	select {
	case <-ctx.Done():
		fmt.Println("timeout ...")
		return errTimeout
	case <-done:
		fmt.Println("exec jobs success")
	}
	return nil
}
