package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var jobs []Job

func init() {
	jobs = make([]Job, 0, 3)

	jobs = append(jobs, func(ctx context.Context) {
		time.Sleep(1 * time.Second)
		fmt.Printf("[ctx: %v] job1 done\n", ctx.Value(ctxCaseKey))
	}, func(ctx context.Context) {
		time.Sleep(3 * time.Second)
		fmt.Printf("[ctx: %v] job2 done\n", ctx.Value(ctxCaseKey))
	}, func(ctx context.Context) {
		time.Sleep(5 * time.Second)
		fmt.Printf("[ctx: %v] job3 done\n", ctx.Value(ctxCaseKey))
	})
}

func Test_ParallelRun(t *testing.T) {
	{
		t.Logf("===== case 1 ====")
		ctx := context.WithValue(context.Background(), ctxCaseKey, "case1")
		ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)
		if err := ParallelRun(ctx, jobs); err != errTimeout {
			t.Logf("ctx 100ms fail, err: %v", err)
			t.FailNow()
		}
		t.Logf("===== case 1 done ====")
	}
	{
		t.Logf("===== case 2 ====")
		ctx := context.WithValue(context.Background(), ctxCaseKey, "case2")
		ctx, _ = context.WithTimeout(ctx, 1*time.Second)
		if err := ParallelRun(ctx, jobs); err != errTimeout {
			t.Logf("ctx 1s fail, err: %v", err)
			t.FailNow()
		}
		t.Logf("===== case 2 done ====")
	}
	{
		t.Logf("===== case 3 ====")
		ctx := context.WithValue(context.Background(), ctxCaseKey, "case3")
		ctx, _ = context.WithTimeout(ctx, 3*time.Second)
		if err := ParallelRun(ctx, jobs); err != errTimeout {
			t.Logf("ctx 3s fail, err: %v", err)
			t.FailNow()
		}
		t.Logf("===== case 3 done ====")
	}
	{
		t.Logf("===== case 4 ====")
		ctx := context.WithValue(context.Background(), ctxCaseKey, "case4")
		ctx, _ = context.WithTimeout(ctx, 6*time.Second)
		if err := ParallelRun(ctx, jobs); err != nil {
			t.Logf("ctx 6s fail, err: %v", err)
			t.FailNow()
		}
		t.Logf("===== case 4 done ====")
	}
}

func Test_ParallelRunWithCancel(t *testing.T) {
	{
		t.Logf("===== case 1 ====")
		ctx := context.WithValue(context.Background(), ctxCaseKey, "case1")
		ctx, cancel := context.WithTimeout(ctx, 6*time.Millisecond)

		go func() {
			// 执行1s就调用下 cancel
			time.Sleep(time.Second)
			cancel()
		}()
		err := ParallelRun(ctx, jobs)
		t.Logf("case 1 err: %v", err)

		t.Logf("===== case 1 done ====")
	}
}
