package search

/*
课程表
https://leetcode.cn/problems/course-schedule/

本质上是，判断是否有环存在
*/
func canFinish(numCourses int, prerequisites [][]int) bool {
	if len(prerequisites) == 0 {
		return true
	}
	vis := make([]bool, numCourses)
	paths := make([][]int, numCourses)
	in := make([]int, numCourses)
	for _, pre := range prerequisites {
		x, y := pre[0], pre[1]
		//if x == y {
		//	continue
		//}
		in[x]++
		paths[y] = append(paths[y], x)
	}
	n := numCourses
	for {
		if n == 0 {
			break
		}
		firsts := make([]int, 0)
		for i, cnt := range in { // 找到所有没有访问过的入度为0的节点
			if vis[i] == false && cnt == 0 {
				firsts = append(firsts, i)
			}
		}
		if len(firsts) == 0 {
			return false
		}
		for _, first := range firsts { // 遍历入度为0的节点，将相邻点的入度减一下
			for _, other := range paths[first] {
				in[other]--
			}
			vis[first] = true
		}
		n -= len(firsts)
	}
	return true
}
