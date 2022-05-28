package concurrent

import (
	"fmt"
	"sync"
	"time"
)

/*
要求按照 12AB34CD... 直到打印到 YZ结束
*/

const (
	sleepTime = 5 * time.Millisecond
)

func PrintNumAndLetter() {
	numCh, letterCh := make(chan struct{}, 0), make(chan struct{}, 0)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go printNum(numCh, letterCh)
	go printLetter(numCh, letterCh, &wg)

	numCh <- struct{}{}
	wg.Wait()
	close(numCh)
	close(letterCh)
}

func printNum(numberCh <-chan struct{}, letterCh chan struct{}) {
	i := 0
	for {
		select {
		case <-numberCh:
			fmt.Printf("%d%d", i, i+1)
			i += 2
			letterCh <- struct{}{}
		}
		time.Sleep(sleepTime)
	}
}

func printLetter(numberCh chan struct{}, letterCh <-chan struct{}, wg *sync.WaitGroup) {
	i := 0
	letterFunc := func(x int) string {
		return string(byte(x + 'A'))
	}
	for {
		select {
		case <-letterCh:
			fmt.Printf("%s%s", letterFunc(i), letterFunc(i+1))
			i += 2
			if i-1 >= 26 {
				fmt.Println()
				wg.Done()
				return
			}
			numberCh <- struct{}{}

		}
		time.Sleep(sleepTime)
	}
}
