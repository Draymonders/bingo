<<<<<<< HEAD
#include <bits/stdc++.h>
=======
#include <vector>
#include <string.h>
>>>>>>> adce5ff... 代码暂存下
using namespace std;

const int N = 11;
string nums;
int dp[N][1024];

void printSpace(int i);

bool isInSet(int v, int mask);
int addToSet(int v, int mask);
int f(int i, int mask, bool isLimit, bool isLeadZero);


/**
如果一个正整数每一个数位都是 <b>互不相同</b> 的，我们称它是 特殊整数 。
给你一个 正整数 n ，请你返回区间 [1, n] 之间特殊整数的数目。
*/
int countSpecialNumbers(int n) {
    nums = to_string(n);
    memset(dp, -1, sizeof(dp));

    return f(0, 0, true, true);
}

/*
f(i, mask, isLimit, isLeadZero):
    i: 位置i
    mask: 到位置i之前所用的数的集合，二进制来表示状态
        i在集合中: (mask>>i)&1
        集合中添加i: mask | (1<<i)
    isLimit: 之前所有访问的前(i-1)位都命中了上界。
        eg: n=124, 如果前面访问了12, 第3位可选的就只有 [0,4] 了
    isLeadZero: 之前所有访问的前(i-1)位都命中了下界。
        eg: n=124, 如果前面访问了00, 第3位0,1,2,3,4 只用统计一次次数
*/
int f(int i, int mask, bool isLimit, bool isLeadZero) {
    if (i == nums.size()) {
        // 题目中，区间计数是按照 [1,n]的，所以全部是0的话，不计数
        return 1 - int(isLeadZero);
    }
    #ifdef debug
        printSpace(i);
        printf("=>");
        printf("i:%d mask:%d isLimit:%d isLeadZero:%d\n", i, mask, int(isLimit), int(isLeadZero));
    #endif
    if (!isLeadZero && !isLimit && dp[i][mask] != -1) {
        return dp[i][mask];
    }

    int ans = 0;
    if (isLeadZero) { // 可以接着前导0
        ans += f(i+1, mask, false, true);
    }
    int up, down;
    up = (isLimit ? nums[i] - '0': 9);
    down = (isLeadZero ? 1 : 0);
    for (int v = down; v <= up; v ++) {
        bool nextLimit = isLimit && (v + '0' == nums[i]); // 下一位是否受限
        if (!isInSet(v, mask)) {
            ans += f(i+1, addToSet(v, mask), nextLimit, false);
        }
    }
    #ifdef debug
        printSpace(i);
        printf("<=");
        printf("i:%d mask:%d up:%d down:%d ans: %d\n", i, mask, up, down, ans);
    #endif
    dp[i][mask] = ans;
    return dp[i][mask];
}

bool isInSet(int v, int mask) {
    return (mask >> v) & 1;
}

int addToSet(int v, int mask) {
    return (mask) | (1 << v);
}


void printSpace(int i) {
    for (int x=0; x<i; x++) {
        printf(" ");
    }
}

int main() {
    vector<int> inputs = vector<int>{20,5,135};
    vector<int> outputs = vector<int>{19,5,110};

    for (int i=0; i<inputs.size(); i++) {
        int input = inputs[i], expect = outputs[i];
        int ans = countSpecialNumbers(input);
        if (ans != expect) {
            printf("input: %d output: %d expect: %d\n", input, ans, expect);
            exit(1);
        }
    }
    printf("success");
    return 0;
}