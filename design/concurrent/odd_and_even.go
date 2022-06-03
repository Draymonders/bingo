package concurrent

import (
	"fmt"
	"sync"
	"time"
)

/*
2个goroutine，打印1...到n
*/

// ParallelPrintOddAndEvenV1
// 版本1：奇数打印完，把chan给偶数，偶数接到chan后，才打印
func ParallelPrintOddAndEvenV1(n int) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	numCh := make(chan struct{})
	go printOdd(n, numCh, &wg)
	go printEven(n, numCh, &wg)

	wg.Wait()
	fmt.Println()
}

// ParallelPrintOddAndEvenV2
// 版本2: 一个buffer chan 存储所有的数字，然后两个 goroutine轮番消费 buffer chan里面的数字
func ParallelPrintOddAndEvenV2(n int) {

	buffers := make(chan int, n)
	numCh := make(chan struct{})

	for i := 1; i <= n; i++ {
		buffers <- i
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go worker1(buffers, numCh, n, wg)
	go worker2(buffers, numCh, n, wg)
	wg.Wait()
	fmt.Println()
}

// printOdd 打印奇数
func printOdd(n int, numCh chan struct{}, wg *sync.WaitGroup) {
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			fmt.Printf("%d", i)
		}
		numCh <- struct{}{}
		time.Sleep(sleepTime)
	}
	wg.Done()
}

// printEven 打印偶数
func printEven(n int, numCh <-chan struct{}, wg *sync.WaitGroup) {
	for i := 1; i <= n; i++ {
		<-numCh
		if i%2 == 0 {
			fmt.Printf("%d", i)
		}
		time.Sleep(sleepTime)
	}
	wg.Done()
}

func worker1(buffers <-chan int, numCh chan struct{}, n int, wg *sync.WaitGroup) {
	for {
		numCh <- struct{}{}
		x := <-buffers
		fmt.Printf("%d", x)

		if x == n {
			wg.Done()
			return
		}
		time.Sleep(sleepTime)
	}
}

func worker2(buffers <-chan int, numCh <-chan struct{}, n int, wg *sync.WaitGroup) {
	for {
		<-numCh
		x := <-buffers
		fmt.Printf("%d", x)
		if x == n {
			wg.Done()
			return
		}
		time.Sleep(sleepTime)
	}
}
