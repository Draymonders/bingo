package tree

import (
	"fmt"

	. "github.com/draymonders/bingo/exam/common"
)

var debug bool

/*
二叉树的最近公共祖先
https://leetcode.cn/problems/lowest-common-ancestor-of-a-binary-tree/
*/
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {

	retNode := &TreeNode{}
	retNodePtr := &retNode
	findLowestParentNode(root, p, q, retNodePtr)
	return *retNodePtr
}

// return val means that root tree (contains root Node & son Nodes) contains node p or q
func findLowestParentNode(root, p, q *TreeNode, retNodePtr **TreeNode) bool {
	defer func() {
		if debug {
			fmt.Printf("root %v retNode %v\n", root, *retNodePtr)
		}
	}()
	if root == nil {
		return false
	}
	l := findLowestParentNode(root.Left, p, q, retNodePtr)
	r := findLowestParentNode(root.Right, p, q, retNodePtr)

	rtEqual := isEqual(root, p) || isEqual(root, q)

	if rtEqual && (l || r) { // 第一种情况，p/q 刚好都是在一条路径上，并且p/q有一个是最近公共祖先
		*retNodePtr = root
		return true
	}
	if l && r { // 第二种情况，p & q 分布在root的两侧
		*retNodePtr = root
		return true
	}
	return rtEqual || l || r
}

func isEqual(a, b *TreeNode) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Val == b.Val
}
