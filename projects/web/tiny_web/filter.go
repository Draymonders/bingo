package main

import (
	"fmt"
	"time"
)

type FilterBuilder func(Filter) Filter

type Filter func(ctx *Context)

func MetricBuilder(next Filter) Filter {
	return func(ctx *Context) {
		st := time.Now()
		next(ctx)
		used := time.Now().Sub(st).Milliseconds()

		fmt.Printf("ctx %+v used %vms\n", ctx, used)
	}
}
