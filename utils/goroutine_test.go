package utils

import (
	"sync"
	"testing"
)

func Test_Recover(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			Recover("recover panic test")
			wg.Done()
		}()
		panic("wow~ go func panic")
	}()

	wg.Wait()

}
