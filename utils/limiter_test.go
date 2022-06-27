package utils

import (
	"testing"
	"time"
)

func Test_Limiter(t *testing.T) {
	l := NewQpsLimiter(time.Second, 1)

	_ = l.Acquire()
	for i := 1; i < 10; i++ {
		if l.Acquire() {
			t.Logf("err: acquire!!")
			t.FailNow()
		}
	}
	time.Sleep(time.Second)
	if !l.Acquire() {
		t.Logf("err: not acquire!!")
		t.FailNow()
	}
}
