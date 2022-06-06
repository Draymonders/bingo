package concurrent

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/draymonders/bingo/utils"
)

// channelClose channel 关闭测试
func channelClose() {

	taskNum := 5
	wg := &sync.WaitGroup{}

	chs := make(chan int, taskNum)
	for i := 1; i <= taskNum; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				utils.Recover(fmt.Sprintf("task-%d", i))
				wg.Done()
			}()

			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

			chs <- i
		}(i)
	}

	go func() {
		// 锁释放的话，释放chs
		wg.Wait()
		close(chs)
	}()

	for ch := range chs {
		fmt.Printf("received %v\n", ch)
	}

	wg.Wait()

}
