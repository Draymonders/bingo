package retry

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	. "github.com/draymonders/bingo/log"

	"github.com/sirupsen/logrus"
)

const (
	retryThreshold = 30
)

var CntZeroErr = errors.New("retry time can't be 0")
var TimeoutZeroErr = errors.New("retry timeout can't be 0")
var TimeoutErr = errors.New("retry timeout")
var ShouldNotRetryErr = errors.New("shouldn't retry")

type runFunc func() error

// 第一个参数表示abort
//    为true, err以后不再继续执行
//    为false, err以后继续重试执行
type runWithAbortFunc func() (bool, error)

func init() {
	Init(logrus.DebugLevel)
}

// Do implement the limit retry logic
func Do(name string, retryCny int, timeout time.Duration, run runFunc) (*stat, error) {
	noAbortRun := func() (bool, error) {
		return false, run()
	}
	return DoCanAbort(name, retryCny, timeout, noAbortRun)
}

func DoCanAbort(name string, retryCnt int, timeout time.Duration, run runWithAbortFunc) (*stat, error) {
	st := &stat{}
	if retryCnt == 0 {
		return st, CntZeroErr
	}
	if timeout == 0 {
		return st, TimeoutZeroErr
	}

	var err error
	timeoutAbort := int32(0) // 超时的标志
	abort := false
	resultCh := make(chan error, 1)
	panicSignal := make(chan string, 1)

	go func() {
		var runErr error
		defer func() {
			panicInfo := recover()
			if panicInfo == nil {
				return
			}
			runErr = errors.New(fmt.Sprint(panicInfo))
			panicSignal <- fmt.Sprintf("[DoCanAbort-Recovery] err=%v debug.Stack:\n%s", runErr, debug.Stack())
		}()
		for i := 0; i < retryCnt; i++ {
			if atomic.LoadInt32(&timeoutAbort) > 0 {
				Log.Infof("stop retry on timeout")
				return
			}

			if i > 0 && !st.shouldRetry() {
				Log.Warnf("Return from error %v, retried %d times", runErr, i)
				resultCh <- ShouldNotRetryErr
				return
			}

			if i > 0 {
				Log.Infof("Retry from error %v, retried %d times", runErr, i)
			}

			// run!!!
			abort, runErr = run()

			if runErr == nil {
				resultCh <- nil
				st.incrSuccess()
				return
			}
			st.incrFail()
			// abort is meaningful only when err is not nil
			if abort {
				resultCh <- runErr
				return
			}
		}
		resultCh <- runErr
		return
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case err = <-resultCh:
		// a read from ch has occurred
	case panicInfo := <-panicSignal:
		// panic!!!
		panic(panicInfo)
	case <-timer.C:
		// the read from ch has timed out
		err = TimeoutErr
		atomic.StoreInt32(&timeoutAbort, 1)
	}

	return st, err
}