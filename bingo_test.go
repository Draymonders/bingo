package bingo

import (
	"fmt"
	"testing"
)

func Test_Panic(t *testing.T) {
	f()
}

func f() {
	defer func() {
		fmt.Println("===呵呵")
	}()
	panic("2333")
}
