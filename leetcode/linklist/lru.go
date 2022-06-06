package linklist

import (
	"fmt"
)

/*
LRU 实现，https://leetcode.cn/problems/lru-cache/
*/
type LRUCache struct {
	len, cap     int
	keyToNode    map[int]*DoubleLinkListNode // key对应的Node
	vHead, vLast *DoubleLinkListNode         // 虚拟节点
}

// 双向链表
type DoubleLinkListNode struct {
	Key       int
	Val       int
	Pre, Next *DoubleLinkListNode
}

func (s *DoubleLinkListNode) String() string {
	return fmt.Sprintf("{key: %d, val: %d}", s.Key, s.Val)
}

func Constructor(capacity int) LRUCache {
	if capacity <= 0 {
		panic("can not generate lruCache")
	}
	vHead, vLast := &DoubleLinkListNode{}, &DoubleLinkListNode{}
	vHead.Next = vLast
	vLast.Pre = vHead
	return LRUCache{
		len:       0,
		cap:       capacity,
		keyToNode: make(map[int]*DoubleLinkListNode),
		vHead:     vHead,
		vLast:     vLast,
	}
}

func (this *LRUCache) Get(key int) int {
	if node, ok := this.keyToNode[key]; ok {
		this.moveToFront(node)
		return node.Val
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if node, ok := this.keyToNode[key]; ok {
		this.moveToFront(node)
		node.Val = value
		return
	}
	if this.len == this.cap {
		lastNodeKey := this.removeLastNode()
		delete(this.keyToNode, lastNodeKey)
	} else {
		this.len++
	}
	node := &DoubleLinkListNode{
		Key: key,
		Val: value,
	}
	this.moveToFront(node)
	this.keyToNode[key] = node
}

func (this *LRUCache) moveToFront(node *DoubleLinkListNode) {
	if node == nil {
		panic("node is nil")
	}
	this.removeNode(node)

	nxt := this.vHead.Next
	this.vHead.Next = node
	node.Pre = this.vHead
	node.Next = nxt
	nxt.Pre = node

}

func (this *LRUCache) removeNode(node *DoubleLinkListNode) {
	if node == nil || node.Pre == nil || node.Next == nil {
		return
	}
	node.Next.Pre = node.Pre
	node.Pre.Next = node.Next
}

func (this *LRUCache) removeLastNode() int {
	if this.len == 0 {
		panic(" can not remove")
	}
	last := this.vLast.Pre
	this.removeNode(last)
	return last.Key
}

func (this *LRUCache) String() string {
	cur := this.vHead.Next
	vals := make([]string, 0)
	for cur != nil {
		vals = append(vals, cur.String())
		cur = cur.Next
		if cur == this.vLast {
			break
		}
	}

	str := fmt.Sprintf("len: %d cap: %d vals: %+v", this.len, this.cap, vals)
	return str
}
