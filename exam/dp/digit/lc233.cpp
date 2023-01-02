#include <bits/stdc++.h>
using namespace std;

const int N = 11;
string nums;
int dp[N][N];
int f(int i, int preOneCnt, bool isLimit, bool isLeadZero);

/**
给定一个整数 n，计算所有小于等于 n 的非负整数中数字 1 出现的个数。
*/
int countDigitOne(int n) {
    nums = to_string(n);
    memset(dp, -1, sizeof(dp));
    return f(0, 0, true, true);
}
/*
f(i, isLimit, isLeadZero):
    i: 位置i
    preOneCnt: preOneCnt 标记前i-1个位置，有1的个数
    isLimit: 之前所有访问的前(i-1)位都命中了上界。
        eg: n=124, 如果前面访问了12, 第3位可选的就只有 [0,4] 了
    isLeadZero: 之前所有访问的前(i-1)位都命中了下界。
        eg: n=124, 如果前面访问了00, 第3位0,1,2,3,4 只用统计一次次数
*/
int f(int i, int preOneCnt, bool isLimit, bool isLeadZero) {
    if (i == nums.size()) {
        if (isLeadZero) {
            return 0;
        }
        return preOneCnt;
    }
    if (!isLimit && !isLeadZero && dp[i][preOneCnt] != -1) {
        return dp[i][preOneCnt];
    }
    int ans = 0;
    if (isLeadZero)
        ans += f(i+1, 0, false, isLeadZero);

    int up, down;
    up = (isLimit) ? nums[i] - '0' : 9;
    down = (isLeadZero) ? 1 : 0;

    for (int v = down; v <= up; v++) {
        int nextOneCnt = (v == 1) ? preOneCnt + 1 : preOneCnt;
        int nextLimit = (isLimit) && (v == up);
        ans += f(i+1, nextOneCnt, nextLimit, false);
    }
    dp[i][preOneCnt] = ans;
    return ans;
}


int main() {
    cout << "LC 233" << endl;
    vector<int> inputs = vector<int>{13, 0};
    vector<int> outputs = vector<int>{6, 0};
    bool allPass = true;

    for (int i=0; i<inputs.size(); i++) {
        int input = inputs[i], expect = outputs[i];
        int ans = countDigitOne(input);
        if (ans != expect) {
            printf("case[%d] input: %d output: %d, but expect: %d\n", i+1, input, ans, expect);
            allPass = false;
            break;
        } else {
            printf("case[%d] success\n", i+1);
        }
    }
    if (allPass) {
        printf("success");
    } else {
        exit(1);
    }

    return 0;
}