package tree

import (
	. "github.com/draymonders/bingo/exam/common"
)

/*
删除二叉搜索树中的节点
https://leetcode.cn/problems/delete-node-in-a-bst/
*/
func deleteNode(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val > key {
		// 删左子树中 节点等于key
		root.Left = deleteNode(root.Left, key)
		return root
	} else if root.Val < key {
		root.Right = deleteNode(root.Right, key)
		return root
	} else {
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}
		// root.Val == key, 找后继节点，替换掉当前根节点
		// eg.  删除key=2的节点
		/*
					  5                               5
			        /    \                           /  \
			       2       7             ====>      3    7
			     / \     /  \                      / \   / \
				1   4   6    10                   1  4  6  10
			       /                                /
				  3                                 3.5
			       \
			       3.5
		*/

		sucNode, sucParentNode := root.Right, root // 后继节点，后继节点的父节点
		for sucNode.Left != nil {
			sucParentNode = sucNode
			sucNode = sucNode.Left
		}

		// 后继节点替换到root节点上
		root.Val = sucNode.Val
		// 将后继节点的右子树接上，删除掉后继节点
		if sucParentNode.Right == sucNode {
			sucParentNode.Right = sucNode.Right
		} else {
			sucParentNode.Left = sucNode.Right
		}
		return root
	}
}

/*
		  5
		 / \
		/   \
       /     \
      2       7
     / \     /  \
	1   4   6    10
       /        /
	  3        8
                \
                 9
*/
func mockBinaryTree() *TreeNode {
	rt := &TreeNode{
		Val: 5,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 1,
			},
			Right: &TreeNode{
				Val: 4,
				Left: &TreeNode{
					Val: 3,
				},
			},
		},
		Right: &TreeNode{
			Val: 7,
			Left: &TreeNode{
				Val: 6,
			},
			Right: &TreeNode{
				Val: 10,
				Left: &TreeNode{
					Val:   8,
					Left:  nil,
					Right: &TreeNode{Val: 9},
				},
			},
		},
	}

	return rt
}

// 中序遍历tree
func inOrderVisit(rt *TreeNode) []int {
	if rt == nil {
		return []int{}
	}
	res := make([]int, 0)
	res = append(res, inOrderVisit(rt.Left)...)
	res = append(res, rt.Val)
	res = append(res, inOrderVisit(rt.Right)...)

	return res
}

/*
给一个target，将小于target的节点都挪到左子树，大于target的节点都挪到右子树，并保持树结构不变
*/
func findNodeAndChangeItToRoot(root *TreeNode, key int) *TreeNode {

	// root.Val == key, 找前驱节点，替换掉当前根节点
	// eg.  key=6
	/*
														 6
														/  \
		              5                               5      7
			        /    \                           /        \
			       2       7             ====>      3         10
			     / \      /  \                      / \
				1   4    6    10                   1  4
			       /      \                           /
				  3       6.5                       3.5
			       \
			       3.5
	*/

	// 解决思路：
	//  找到target节点，并找到前驱节点，以前驱节点构建左子树
	//      同理，找到后继节点，以后继节点构建右子树
	//  前驱节点和后继节点分别作为targetNode的左儿子和右儿子

}
