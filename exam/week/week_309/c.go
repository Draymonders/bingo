package main

import (
	"fmt"
	"os"
)

const N = 31

/**
lc 6169:
给你一个由 正 整数组成的数组 nums 。

如果nums 的子数组中位于 不同 位置的每对元素按位 与（AND）运算的结果等于 0 ，则称该子数组为 优雅 子数组。

返回 最长 的优雅子数组的长度。

子数组 是数组中的一个 连续 部分。

注意：长度为 1 的子数组始终视作优雅子数组。
*/
func longestNiceSubarray(nums []int) int {
	l, r, n := 0, 0, len(nums)
	ans := 0
	bits := make([]int, N+1)

	// [l, r]
	for r < n {
		for i := 0; i < N; i++ {
			if ((nums[r] >> i) & 1) > 0 {
				bits[i]++
			}
		}
		for l <= r && !check(bits) {
			for i := 0; i < N; i++ {
				if ((nums[l] >> i) & 1) > 0 {
					bits[i]--
				}
			}
			l++
		}
		if r-l+1 > ans {
			ans = r - l + 1
		}
		r++
	}
	return ans
}

func check(bits []int) bool {
	for i := 0; i < N; i++ {
		if bits[i] > 1 {
			return false
		}
	}
	return true
}

func main() {
	type arg struct {
		Nums []int
		Want int
	}

	args := []*arg{{Nums: []int{1, 3, 8, 48, 10}, Want: 3}, {Nums: []int{3, 1, 5, 11, 13}, Want: 1}}

	for _, tmp := range args {
		if v := longestNiceSubarray(tmp.Nums); v != tmp.Want {
			os.Exit(1)
		}
	}
	fmt.Println("success!")
}
