package concurrent

import (
	"fmt"
	"sync"
	"time"
)

/*
m个goroutine 打印 1...N
*/

// ParallelPrintOntToN m个goroutine 打印 1...n
func ParallelPrintOntToN(m, n int) {
	if m > n {
		m = n
	}
	// 定义m把锁，第i个goroutine拥有第i把锁，锁用来控制第(i+1)个goroutine的写
	numChs := make([]chan struct{}, 0, m)
	for i := 1; i <= m; i++ {
		numChs = append(numChs, make(chan struct{}))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < m; i++ {
		go printNumber(i, m, n, numChs, wg)
	}
	// 先给第1个goroutine释放掉锁
	numChs[m-1] <- struct{}{}

	wg.Wait()
	fmt.Println()
}

// printNum idx表示第i个worker，0 <= idx < m
func printNumber(idx, m, n int, numChs []chan struct{}, wg *sync.WaitGroup) {
	pre := func(i int) int {
		return ((i-1)%m + m) % m
	}
	cur := func(i int) int {
		return (i%m + m) % m
	}
	// fmt.Printf("idx: %d pre(idx): %d cur(idx): %d\n", idx, pre(idx), cur(idx))

	for i := 1; i <= n; i++ {
		if (i-1)%m == idx {
			// 获取锁
			<-numChs[pre(idx)]
			fmt.Printf("%d", i)
			if i == n {
				wg.Done()
				return
			}
			numChs[cur(idx)] <- struct{}{}
		}
		time.Sleep(sleepTime)
	}
}
