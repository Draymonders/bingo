package retry

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestStat(t *testing.T) {
	st := &stat{}
	st.incrSuccess()
	if st.getSuccess() != 1 {
		t.Fatal("success must be 1")
	}
	t.Logf("%v", st)
	st.incrFail()
	if st.getFail() != 1 {
		t.Fatal("fail must be 1")
	}
	t.Logf("%v", st)
}

// 重试参数设置为0
func TestRetryCntZeroError(t *testing.T) {
	_, err := Do("TestRetryTimeout", 0, time.Second, func() error {
		return nil
	})
	if err != CntZeroErr {
		t.Fatal("Error should be CntZeroErr")
	}
	t.Logf("CntZeroErr %v", err)
}

// 超时参数设置为0
func TestRetryTimeoutZeroError(t *testing.T) {
	_, err := Do("TestRetryTimeout", 3, 0, func() error {
		return nil
	})
	if err != TimeoutZeroErr {
		t.Fatal("Error should be TimeoutZeroErr")
	}
	t.Logf("TimeoutZeroErr %v", err)
}

// 方法超时，不可abort
func TestRetryTimeout(t *testing.T) {
	stat, err := Do("TestRetryTimeout", 3, time.Second, func() error {
		fmt.Println("sleep 2s")
		time.Sleep(time.Second * 2)
		return nil
	})
	if err != TimeoutErr {
		t.Fatal("Error should be TimeoutErr")
	}
	t.Logf("TimeoutErr %v, %v", err, stat)
}

// 方法超时，可abort
func TestRetryCanAbortTimeout(t *testing.T) {
	stat, err := DoCanAbort("TestRetryCanAbortTimeout", 3, time.Second, func() (bool, error) {
		fmt.Println("sleep 2s")
		time.Sleep(time.Second * 2)
		return true, nil
	})
	if err != TimeoutErr {
		t.Fatal("Error should be TimeoutErr")
	}
	t.Logf("TimeoutErr %v, %v", TimeoutErr, stat)
}

// 执行成功，可abort
func TestRetryCanAbortSuccess(t *testing.T) {
	stat, err := DoCanAbort("TestRetryCanAbortSuccess", 3, time.Second, func() (bool, error) {
		return true, nil
	})
	if err != nil {
		t.Fatal("Error should be nil")
	}
	t.Logf("err is %v, %v", err, stat)
}

// 执行成功，不可abort
func TestRetrySuccess(t *testing.T) {
	stat, err := Do("TestRetrySuccess", 3, time.Second, func() error {
		return nil
	})
	if err != nil {
		t.Fatal("Error should be nil")
	}
	t.Logf("err is %v, %v", err, stat)
}

// 方法执行错误，可abort
func TestRetryCanAbortFailed(t *testing.T) {
	funcErr := errors.New("funcErr")
	stat, err := DoCanAbort("TestRetryCanAbortFailed", 3, time.Second*10, func() (bool, error) {
		return true, funcErr
	})
	if err != funcErr {
		t.Fatalf("Error should be equal to funcErr, but actually is %v", err)
	}
	t.Logf("err is %v, %v", err, stat)
	// err is funcErr, stat successCnt: 0, failCnt: 1
}

// 方法执行错误，可abort
func TestRetryFailed(t *testing.T) {
	funcErr := errors.New("funcErr")
	stat, err := Do("TestRetryFailed", 3, time.Second*10, func() error {
		return funcErr
	})
	if err != funcErr {
		t.Fatalf("Error should be funcErr, but actually is %v", err)
	}
	t.Logf("err is %v, %v", err, stat)
	// err is funcErr, stat successCnt: 0, failCnt: 3
}

// 超时情况下，执行了多少次
func TestStopRetryWhenTimeout(t *testing.T) {
	var count uint32
	stat, err := DoCanAbort("TestStopRetryWhenTimeout", 5, time.Second/2, func() (bool, error) {
		atomic.AddUint32(&count, 1) // count ++

		time.Sleep(time.Second)
		return false, errors.New("some err")
	})

	// wait 3 seconds to see if retry stop
	time.Sleep(time.Second * 3)
	if err == nil {
		t.Errorf("Do error nil, expect retry timeout")
	}

	// should not retry on timeout
	if atomic.LoadUint32(&count) > 1 {
		t.Errorf("count expect: %v, got: %v", 1, count)
	}
	t.Logf("err is %v, %v", err, stat)
	// err is retry timeout, stat successCnt: 0, failCnt: 1
}
