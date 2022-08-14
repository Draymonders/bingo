#include <bits/stdc++.h>
using namespace std;


const int N = 35;
string nums;
int dp[N][N];
int f(int i, int mask, bool isLimit, bool isLeadZero);

/**
给定一个正整数 n ，返回范围在 [0, n] 都非负整数中，其二进制表示不包含 连续的 1 的个数。

====> 枚举二进制！！！
*/
int findIntegers(int n) {
    int m = n;
    nums.clear();
    while (m > 0) {
        nums.push_back(char('0' + (m%2)));
        m /= 2;
    }
    reverse(nums.begin(), nums.end());
    memset(dp, -1, sizeof(dp));

    #ifdef debug
        cout << "nums: " << nums << endl;
    #endif

    return f(0, 0, true, true);
}

void printSpace(int i) {
    for (int x=0; x<i; x++) {
        printf(" ");
    }
}

/*
f(i, mask, isLimit, isLeadZero):
    i: 位置i
    mask: i-1位置的填的数
    isLimit: 之前所有访问的前(i-1)位都命中了上界。
        eg: n=124, 如果前面访问了12, 第3位可选的就只有 [0,4] 了
    isLeadZero: 之前所有访问的前(i-1)位都命中了下界。
        eg: n=124, 如果前面访问了00, 第3位0,1,2,3,4 只用统计一次次数
*/
int f(int i, int mask, bool isLimit, bool isLeadZero) {
    if (i == nums.size()) {
        return 1;
    }
    #ifdef debug
        printSpace(i);
        printf("=>");
        printf("i:%d mask:%d isLimit:%d isLeadZero:%d\n", i, mask, int(isLimit), int(isLeadZero));
    #endif

    if (!isLimit && !isLeadZero && dp[i][mask] != -1) {
        return dp[i][mask];
    }

    int ans = 0;
    if (isLeadZero)
        ans += f(i+1, 0, false, isLeadZero);
    #ifdef debug
        printSpace(i);
        printf("<==");
        printf("i:%d mask:%d ans: %d\n", i, mask, ans);
    #endif
    int up, down;
    up = (isLimit) ? nums[i] - '0' : 1;
    down = (isLeadZero) ? 1 : 0;

    for (int v = down; v <= up; v++) {
        int nextMask = v;
        int nextLimit = (isLimit) && (v == up);
        if (mask && v == 1) {
            continue;
        }
        ans += f(i+1, nextMask, nextLimit, false);
    }
    dp[i][mask] = ans;
    #ifdef debug
        printSpace(i);
        printf("<=");
        printf("i:%d mask:%d up:%d down:%d ans: %d\n", i, mask, up, down, ans);
    #endif
    return ans;
}

int main() {
    cout << "LC 600" << endl;
    vector<int> inputs=vector<int>{5,1,2,28};
    vector<int> outputs = vector<int>{5,2,3,13};
    bool allPass = true;

    for (int i=0; i<inputs.size(); i++) {
        int input = inputs[i], expect = outputs[i];
        int ans = findIntegers(input);
        if (ans != expect) {
            printf("case[%d] input: %d output: %d, but expect: %d\n", i, input, ans, expect);
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