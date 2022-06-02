package sort

import "fmt"

var debug bool

/*
数组中的第K个最大元素
https://leetcode.cn/problems/kth-largest-element-in-an-array/
*/
func findKthLargest(nums []int, k int) int {
	if debug {
		fmt.Println(nums)
	}
	return quickSort(nums, 0, len(nums)-1, len(nums)-k)
}

// nums 区间 [l, r] 找到 位置为k的数
func quickSort(nums []int, l int, r int, k int) int {
	if l == r {
		if k != l {
			panic("l == r && k != l")
		}
		return nums[l]
	}
	i, j := l, r
	tmpVal := nums[l]
	for i < j {
		for i < j && nums[j] >= tmpVal {
			j--
		}
		for i < j && nums[i] <= tmpVal {
			i++
		}
		if i < j {
			nums[i], nums[j] = nums[j], nums[i]
		}
	}

	// make sure: i == j && nums[l] >= nums[i], so swap(nums[l], nums[i]) is ok
	if nums[l] < nums[i] {
		panic(fmt.Sprintf("nums: %v l %d r %d i %d", nums, l, r, i))
	}
	nums[i], nums[l] = nums[l], nums[i]

	if debug {
		fmt.Printf("-> l: %d r: %d i: %d j: %d \n", l, r, i, j)
		for a := l; a <= r; a++ {
			if a == i {
				fmt.Printf(" |%d|", nums[a])
			} else {
				fmt.Printf(" %d ", nums[a])
			}
		}
		fmt.Println()
	}

	// make sure: nums[l,l+1...i-1] <= nums[i] <= nums[i+1,i+2..r]
	if i == k {
		return nums[i]
	} else if i > k {
		return quickSort(nums, l, i-1, k)
	}
	return quickSort(nums, i+1, r, k)
}
