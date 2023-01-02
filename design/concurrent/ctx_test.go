package concurrent

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCtxCancel(t *testing.T) {
	DoCtxCancel()
}

func Test_ChanDone(t *testing.T) {
	errCh := make(chan error, 2)
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case err, ok := <-errCh:
				fmt.Println("ok", ok, "err", err)
				if !ok {
					return
				}
			}
		}
	}()
	errCh <- errors.New("mock err")
	go func() {
		wg.Add(1)
		defer wg.Done()
		close(errCh)
	}()

	wg.Wait()
	time.Sleep(50 * time.Millisecond)
}
