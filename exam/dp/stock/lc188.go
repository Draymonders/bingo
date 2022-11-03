package main

import (
	"fmt"
	"os"
)

/**
给定一个整数数组prices ，它的第 i 个元素prices[i] 是一支给定的股票在第 i 天的价格。

设计一个算法来计算你所能获取的最大利润。你最多可以完成 k 笔交易。

注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。
*/

const N = 1010
const M = 110

func maxProfit(k int, prices []int) int {
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			for k := 0; k < 2; k++ {
				dp[i][j][k] = -1
			}
		}
	}
	return recur(prices, k, 0, len(prices), 0, notOwn)
}

// 记忆化
var dp [N][M][2]int // dp[st][curCnt][chooseType]

type chooseType int

const (
	own chooseType = iota
	notOwn
)

/**

choose 1: 当前持有股票；0： 当前未持有股票
curCnt: 当前交易了几笔
*/
func recur(prices []int, k, st, ed, curCnt int, choose chooseType) int {
	if st == ed {
		return 0
	}
	if dp[st][curCnt][choose] != -1 {
		return dp[st][curCnt][choose]
	}
	var v1, v2 int
	if choose == own {
		v1 = recur(prices, k, st+1, ed, curCnt, notOwn) + prices[st] // 卖掉
		v2 = recur(prices, k, st+1, ed, curCnt, own)
	} else {
		if curCnt < k {
			v1 = recur(prices, k, st+1, ed, curCnt+1, own) - prices[st]
		}
		v2 = recur(prices, k, st+1, ed, curCnt, notOwn)
	}
	dp[st][curCnt][choose] = maxInt(v1, v2)
	return dp[st][curCnt][choose]
}

func maxInt(v1, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}

func main() {
	type Case struct {
		K      int
		Prices []int
		Want   int
	}

	cases := []*Case{
		{K: 2, Prices: []int{2, 4, 1}, Want: 2},
		{K: 2, Prices: []int{3, 2, 6, 5, 0, 3}, Want: 7},
	}
	for i, caze := range cases {
		real := maxProfit(caze.K, caze.Prices)
		if real != caze.Want {
			fmt.Printf("case %d: want: %d but real: %d\n", i, caze.Want, real)
			os.Exit(1)
		}
	}
	fmt.Println("Success!!!")
}
