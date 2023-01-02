package string

/*
无重复字符的最长子串
https://leetcode.cn/problems/longest-substring-without-repeating-characters/
*/
func lengthOfLongestSubstring(s string) int {
	idx := make(map[int]int, 0)

	res := 0
	pre := -1 // 表示当前 [pre+1, i]是个合法区间
	for i := 0; i < len(s); i++ {
		ch := int(s[i])
		if p, ok := idx[ch]; ok {
			pre = max(pre, p)
		}
		idx[ch] = i
		res = max(res, i-pre)
	}

	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
