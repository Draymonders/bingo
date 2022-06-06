package array

func maxSubArray(nums []int) int {
	res := -0x3f3f3f3f
	ans := 0
	for i := 0; i < len(nums); i++ {
		res = max(res, nums[i])
		if ans < 0 {
			ans = 0
		}
		ans += nums[i]

		res = max(res, ans)
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
