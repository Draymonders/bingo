package utils

import (
	"runtime/debug"

	. "github.com/draymonders/bingo/log"
)

func Recover(name string) {
	if r := recover(); r != nil {
		// metrics
		Log.Errorf("[panic] Recover capture, name: %v, recover: %v\n stack: \n%v\n", name, r, string(debug.Stack()))
	}
}
