package observer

import (
	"strconv"
	"sync"
	"testing"
)

func Test_Notifier(t *testing.T) {
	n := 5

	var shopNotifiers []*ShopNotifier
	for i := 1; i <= n; i++ {
		shopTypeVal := i % 2
		shop := &Shop{
			Id:       int64(i + 10000),
			Name:     "店铺" + strconv.FormatInt(int64(i), 10) + "号",
			ShopType: ShopType(shopTypeVal),
		}

		shopNotifiers = append(shopNotifiers, getShopNotifier(shop))
	}

	wg := &sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(notifier *ShopNotifier) {
			_ = notifier.NotifyAll()
			wg.Done()
		}(shopNotifiers[i])
	}
	wg.Wait()
}
