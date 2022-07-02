package group

import (
	"fmt"
	"time"
)

// 场景：运行多个job，有一个执行完的情况下，就停止所有
func waitOne() {
	done := make(chan error, 2)
	stop := make(chan struct{})
	go func() {
		done <- serveHTTP(stop)
	}()
	go func() {
		done <- serveGRPC(stop)
	}()

	// 两种情况，随便来了个done
	for i := 0; i < cap(done); i++ {
		<-done
		if i == 0 {
			// 关闭 channel 会产生一个广播机制，所有向 channel 读取消息的 goroutine 都会收到消息。
			close(stop)
		}
	}
}

func serveHTTP(stop <-chan struct{}) error {
	done := make(chan struct{}, 1)
	go func() {
		fmt.Println("start serve http")
		time.Sleep(10 * time.Second)

		done <- struct{}{}
	}()

	select {
	case <-stop:
		fmt.Println("http serve stop, due to stop chan")
	case <-done:
		fmt.Println("stop serve http")
	}
	return nil
}

func serveGRPC(stop <-chan struct{}) error {
	done := make(chan struct{}, 1)
	go func() {
		fmt.Println("start serve grpc")
		time.Sleep(1 * time.Second)
		done <- struct{}{}
	}()

	select {
	case <-stop:
		fmt.Println("grpc serve stop, due to stop chan")
	case <-done:
		fmt.Println("stop serve grpc")
	}
	return nil
}
