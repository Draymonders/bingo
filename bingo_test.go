package bingo

import (
	"encoding/json"
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

// bool 单个值也可以反序列化
func Test_JsonWithBool(t *testing.T) {
	{
		val := "true"
		var res bool
		err := json.Unmarshal([]byte(val), &res)
		if err != nil {
			t.FailNow()
		}
		t.Log("bool", res)
	}

	{
		val := "1"
		var v int
		if err := json.Unmarshal([]byte(val), &v); err != nil {
			t.Fatalf("val: %v, unmarshal err: %v", val, err)
		}
		t.Log("int", v)
	}

}
