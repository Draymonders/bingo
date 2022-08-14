package sort

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/draymonders/bingo/exam/common"
)

func Test_MultiTopK(t *testing.T) {
	debug = false

	type args struct {
		nums []int
		k    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		//{
		//	name: "case 1",
		//	args: args{
		//		nums: []int{3, 2, 1, 5, 6, 4},
		//		k:    2,
		//	},
		//	want: 5,
		//},
		//{
		//	name: "case 2",
		//	args: args{
		//		nums: []int{3, 2, 3, 1, 2, 4, 5, 5, 6},
		//		k:    4,
		//	},
		//	want: 4,
		//},
	}

	n, sliceSize, maxVal := 1000, 100, 100
	for i := 0; i < n; i++ {
		nums := common.GenerateIntSlice(sliceSize, maxVal)
		k := rand.Intn(sliceSize)
		for k == 0 {
			k = rand.Intn(sliceSize)
		}
		tests = append(tests, struct {
			name string
			args args
			want int
		}{
			name: "case " + strconv.FormatInt(int64(i), 10), args: args{
				nums: nums,
				k:    k,
			}, want: getKth(nums, k)})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("name: %v arg: %v", tt.name, tt.args)
			got := findKthLargest(tt.args.nums, tt.args.k)
			if got != tt.want {
				t.Errorf("findKthLargest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getKth(nums []int, k int) int {
	res := make([]int, len(nums))
	copy(res, nums)

	sort.Ints(res)
	fmt.Printf("origin: %v res: %v k: %d retVal: %d\n", nums, res, k, res[len(res)-k])

	return res[len(res)-k]
}

func Test_SingleTopK(t *testing.T) {
	debug = true
	nums := []int{11, 62, 89, 28, 74, 11, 45, 37, 6, 95}
	k := 6
	//nums := []int{81, 87, 47, 59, 81, 18, 25, 40, 56, 0}
	//k := 4
	v := findKthLargest(nums, k)
	fmt.Printf("res: %d\n", v)

	getKth(nums, k)
}
