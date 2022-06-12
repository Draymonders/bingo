package single_flight

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

/*
singleFlight 对相同的key，进行了队列 + waitGroup 的管理。waitGroup.Wait() 完成后，将队列里所有的请求都进行返回
*/

// cachePenetration 避免缓存穿透
func cachePenetration() {
	group := &singleflight.Group{}
	cnts := &sync.Map{}
	wg := &sync.WaitGroup{}
	const n = 5
	//chs := make(chan struct{}, n)
	//defer close(chs)

	// n个协程去拿key，看看最终计算了多少次
	for i := 1; i <= n; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)

			key := fmt.Sprintf("k_%d", i%2)
			v, err, share := group.Do(key, func() (interface{}, error) {
				vv := 1
				if v, ok := cnts.Load(key); ok {
					vv = v.(int) + 1
				}
				cnts.Store(key, vv)
				return key, nil
			})
			fmt.Printf("%d key %v, v %v err %v share %v\n", i, key, v, err, share)
		}(i)
	}
	wg.Wait()

	cnts.Range(func(key, value interface{}) bool {
		fmt.Printf("cnts: key %v val: %v\n", key, value)
		return true
	})

}
