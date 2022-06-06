package linklist

import (
	"testing"

	"github.com/draymonders/bingo/leetcode/common"
)

// 翻转链表测试
func Test_reverseList(t *testing.T) {
	st, ed := 1, 5
	var head *common.ListNode
	var cur *common.ListNode
	for i := st; i < ed; i++ {
		if i == st {
			cur = &common.ListNode{Val: i, Next: nil}
			head = cur
		} else {
			cur.Next = &common.ListNode{Val: i, Next: nil}
			cur = cur.Next
		}
	}

	debug = true
	head = reverseList(head)

	PrintListNode("reverseListResult", head)
}

// LRU测试
func Test_Lru(t *testing.T) {

	/*
		{
				cache := Constructor(2)

				cache.Put(1, 1)
				t.Logf("cache: %v", cache.String())
				cache.Put(2, 2)
				t.Logf("cache: %v", cache.String())

				cache.Put(1, 3)
				t.Logf("cache: %v", cache.String())

				cache.Put(3, 3)
				t.Logf("cache: %v", cache.String())
			}
	*/

	{
		// bad case
		//["LRUCache","put","get","put","get","get"]
		//[[1],[2,1],[2],[3,2],[2],[3]]
		cache := Constructor(1)
		cache.Put(2, 1)
		cache.Get(2)
		t.Logf("cache: %v", cache.String())
		cache.Put(3, 2)
		t.Logf("cache: %v", cache.String())
		cache.Get(2)
		cache.Get(3)
		t.Logf("cache: %v", cache.String())

	}
}
