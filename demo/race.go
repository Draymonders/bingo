package main

import (
	"fmt"
	"time"
)

func request(c *string) {
	println(fmt.Sprintf("fullPath: %s", *c))
}

func main() {

	fullPath := "init"

	go func() {
		for {
			request(&fullPath)
		}
	}()

	for {
		fullPath = ""
		time.Sleep(10 * time.Nanosecond)
		fullPath = "/test/test/test"
		time.Sleep(10 * time.Nanosecond)
	}

}
