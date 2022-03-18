package retry

import (
	"fmt"
	"sync/atomic"
)

type stat struct {
	success atomic.Value
	fail    atomic.Value
}

func (s *stat) shouldRetry() bool {
	// todo 统计时间窗口内的failTimes, successTimes
	failedTimes := s.getFail()
	if failedTimes > retryThreshold && failedTimes*10 > s.getSuccess() {
		return false
	}
	return true
}

func (s *stat) incrSuccess() {
	s.success.Store(s.getSuccess() + 1)
}

func (s *stat) incrFail() {
	s.fail.Store(s.getFail() + 1)
}

func (s *stat) getSuccess() int64 {
	if v := s.success.Load(); v != nil {
		return v.(int64)
	}
	return 0
}

func (s *stat) getFail() int64 {
	if v := s.fail.Load(); v != nil {
		return v.(int64)
	}
	return 0
}

func (s *stat) String() string {
	return fmt.Sprintf("stat successCnt: %d, failCnt: %d", s.getSuccess(), s.getFail())
}
