package tree

import (
	. "github.com/draymonders/bingo/exam/common"
)

/*
二叉树中的最大路径和
https://leetcode.cn/problems/binary-tree-maximum-path-sum/
*/

func maxPathSum(root *TreeNode) int {
	mx := -0x3f3f3f3f

	dfs(root, &mx)
	return mx
}

// returnVal is the max(max(leftSonPathSum. rightSonPathSum), 0) + root.Val
func dfs(root *TreeNode, mxPtr *int) int {
	if root == nil {
		return 0
	}
	leftSum, rightSum := dfs(root.Left, mxPtr), dfs(root.Right, mxPtr)

	*mxPtr = max(*mxPtr, root.Val)
	*mxPtr = max(*mxPtr, root.Val+leftSum)
	*mxPtr = max(*mxPtr, root.Val+rightSum)
	*mxPtr = max(*mxPtr, root.Val+leftSum+rightSum)

	return max(max(leftSum, rightSum), 0) + root.Val
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
