package tree

import (
	"reflect"
	"testing"

	. "github.com/draymonders/bingo/exam/common"
)

// 测试最大路径和
func Test_maxPathSum(t *testing.T) {
	{
		rt := &TreeNode{
			Val: 1,
			Left: &TreeNode{
				Val:   2,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{
				Val:   3,
				Left:  nil,
				Right: nil,
			},
		}
		retVal := maxPathSum(rt)
		want := 6

		if retVal != want {
			t.Logf("retVal %d want %d", retVal, want)
			t.FailNow()
		}
	}

	{
		rt := &TreeNode{
			Val: -10,
			Left: &TreeNode{
				Val:   9,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{
				Val: 20,
				Left: &TreeNode{
					Val:   15,
					Left:  nil,
					Right: nil,
				},
				Right: &TreeNode{
					Val:   7,
					Left:  nil,
					Right: nil,
				},
			},
		}
		retVal := maxPathSum(rt)
		want := 42

		if retVal != want {
			t.Logf("retVal %d want %d", retVal, want)
			t.FailNow()
		}
	}
}

func Test_findLowestParentNode(t *testing.T) {
	rt := &TreeNode{
		Val: -10,
		Left: &TreeNode{
			Val:   9,
			Left:  nil,
			Right: nil,
		},
		Right: &TreeNode{
			Val: 20,
			Left: &TreeNode{
				Val:   15,
				Left:  nil,
				Right: nil,
			},
			Right: &TreeNode{
				Val:   7,
				Left:  nil,
				Right: nil,
			},
		},
	}
	debug = true

	{
		pVal := 7
		qVal := 15
		want := 20
		node := lowestCommonAncestor(rt, &TreeNode{Val: pVal}, &TreeNode{Val: qVal})
		if node == nil || node.Val != want {
			t.Logf("node %d & node %d is %d, but node %v", pVal, qVal, want, node)
			t.FailNow()
		}
	}
	//{
	//	pVal := 20
	//	qVal := 15
	//	want := 20
	//	node := lowestCommonAncestor(rt, &TreeNode{Val: pVal}, &TreeNode{Val: qVal})
	//	if node == nil || node.Val != want {
	//		t.Logf("node %d & node %d is %d, but node %v", pVal, qVal, want, node)
	//		t.FailNow()
	//	}
	//}
	//
	//{
	//	pVal := 9
	//	qVal := 7
	//	want := -10
	//	node := lowestCommonAncestor(rt, &TreeNode{Val: pVal}, &TreeNode{Val: qVal})
	//	if node == nil || node.Val != want {
	//		t.Logf("node %d & node %d is %d, but node %v", pVal, qVal, want, node)
	//		t.FailNow()
	//	}
	//}
}

// 二叉搜索树删除某一节点
func Test_delNode(t *testing.T) {

	removeValF := func(lst []int, val int) []int {
		newLst := make([]int, 0)
		for _, v := range lst {
			if v == val {
				continue
			}
			newLst = append(newLst, v)
		}
		return newLst
	}
	//
	//{
	//	rt := mockBinaryTree()
	//	removeVal := 5
	//	lst := inOrderVisit(rt)
	//	fmt.Println("mockList", lst)
	//	want := removeValF(lst, removeVal)
	//	fmt.Printf("removeVal %d want %v\n", removeVal, want)
	//	rt = deleteNode(rt, removeVal)
	//	got := inOrderVisit(rt)
	//
	//	if !reflect.DeepEqual(got, want) {
	//		fmt.Printf("got %v\n", got)
	//		t.FailNow()
	//	}
	//	fmt.Println("After mockList", lst)
	//
	//}
	//fmt.Println("===========")
	//{
	//	rt := mockBinaryTree()
	//	removeVal := 6
	//	lst := inOrderVisit(rt)
	//	want := removeValF(lst, removeVal)
	//	rt = deleteNode(rt, removeVal)
	//	got := inOrderVisit(rt)
	//	if !reflect.DeepEqual(got, want) {
	//		t.Logf("want %v but got %v", want, got)
	//		t.FailNow()
	//	}
	//	fmt.Println("After mockList", lst)
	//
	//}
	//debug = true
	for i := 1; i <= 10; i++ {
		t.Logf("======")
		t.Logf("removeVal %v start", i)
		rt := mockBinaryTree()
		removeVal := i
		lst := inOrderVisit(rt)
		want := removeValF(lst, removeVal)
		rt = deleteNode(rt, removeVal)
		got := inOrderVisit(rt)
		if !reflect.DeepEqual(got, want) {
			t.Logf("removeVal %v want %v but got %v", removeVal, want, got)
			t.FailNow()
		}
		t.Logf("======")
	}
}
