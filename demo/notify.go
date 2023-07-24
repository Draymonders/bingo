package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stop := make(chan struct{}, 1)
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		stop <- struct{}{}
	}()

	<-stop
	fmt.Println("2333")
}
