package loader

import (
	"fmt"
	"strconv"
	"sync"
)

type item struct {
	Id int64 // 资质Id

	AThing string // A的信息
	BThing string // B的信息
}

type ILoader interface {
	Init(ids []int64)        // load初始化
	Load()                   // 加载数据
	Fill(item []interface{}) // 加载完数据后，填充到
}

type ALoader struct {
	ids []int64
	M   map[int64]string
}

func (s *ALoader) Init(ids []int64) {
	s.ids = ids
	s.M = make(map[int64]string, 0)
}

func (s *ALoader) Load() {
	for _, id := range s.ids {
		s.M[id] = "A_" + strconv.FormatInt(id, 10)
	}
}

func (s *ALoader) Fill(items []interface{}) {
	if len(items) == 0 {
		return
	}
	for _, itm := range items {
		if it, ok := itm.(*item); ok {
			if str, ok2 := s.M[it.Id]; ok2 {
				it.AThing = str
			}
		}
	}
}

type BLoader struct {
	ids []int64
	M   map[int64]string
}

func (s *BLoader) Init(ids []int64) {
	s.ids = ids
	s.M = make(map[int64]string, 0)
}

func (s *BLoader) Load() {
	for _, id := range s.ids {
		s.M[id] = "B_" + strconv.FormatInt(id*15, 10)
	}
}

func (s *BLoader) Fill(items []interface{}) {
	if len(items) == 0 {
		return
	}
	for _, itm := range items {
		if it, ok := itm.(*item); ok {
			if str, ok2 := s.M[it.Id]; ok2 {
				it.BThing = str
			}
		}
	}
}

type Worker interface {
	Run(item interface{}) // 执行item
}

// LoadAndExec 不好的地方就是还得外面强转一下
func LoadAndExec(ids []int64, loaders []ILoader) []interface{} {
	//loaders := []ILoader{&ALoader{}, &BLoader{}}
	wg := sync.WaitGroup{}
	//ids := []int64{233, 244}
	var items []interface{}
	for i := 0; i < len(ids); i++ {
		items = append(items, &item{Id: ids[i]})
	}

	for _, loader := range loaders {
		loader.Init(ids)

	}
	for _, loader := range loaders {
		wg.Add(1)
		go func(ld ILoader) {
			defer wg.Done()
			ld.Load()
		}(loader)
	}
	wg.Wait()

	for _, loader := range loaders {
		loader.Fill(items)
	}

	for _, itm := range items {
		fmt.Printf("item: %+v\n", itm)
	}
	return items
}
