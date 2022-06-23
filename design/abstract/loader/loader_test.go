package loader

import (
	"fmt"
	"testing"

	"github.com/draymonders/bingo/utils"
)

func Test_LoadAndExec(t *testing.T) {
	itemItfs := LoadAndExec([]int64{111, 222}, []ILoader{&ALoader{}, &BLoader{}})

	var items []*item
	for _, itm := range itemItfs {
		if it, ok := itm.(*item); ok {
			items = append(items, it)
		}
	}
	fmt.Println(utils.ObjToJson(items))
}
