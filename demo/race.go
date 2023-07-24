package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	doCtx()
}

func doCtx() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, "k1", "v1")
	ctx = context.WithValue(ctx, "k1", "v2")

	v := ctx.Value("k1").(string)
	fmt.Println(v)
}

// 竞争，会panic
func race() {
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

func request(c *string) {
	println(fmt.Sprintf("fullPath: %s", *c))
}
