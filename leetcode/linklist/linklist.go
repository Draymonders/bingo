package linklist

import (
	"fmt"

	. "github.com/draymonders/bingo/leetcode/common"
)

var debug bool

/*
反转链表
https://leetcode.cn/problems/reverse-linked-list/
*/
func reverseList(head *ListNode) *ListNode {
	var pre, cur *ListNode
	cur = head

	for cur != nil {
		next := cur.Next

		cur.Next = pre
		pre = cur
		cur = next
		if debug {
			PrintListNode("pre", pre)
			PrintListNode("cur", cur)
		}
	}
	return pre
}

func PrintListNode(logPre string, cur *ListNode) {
	fmt.Printf("%s: ", logPre)
	for cur != nil {
		if cur.Next == nil {
			fmt.Printf("%d", cur.Val)
		} else {
			fmt.Printf("%d->", cur.Val)
		}

		cur = cur.Next
	}
	fmt.Println()
}

// 1 -> 2 -> 3 -> 4
// begin:
//    pre: nil, cur: 1->2->3->4
// first:
//    pre: 1, cur:2->3->4
// second:
//    pre: 2->1, cur: 3->4
// third:
//    pre: 3->2->1, cur: 4
// fourth:
//    pre: 4->3->2->1, cur: nil
