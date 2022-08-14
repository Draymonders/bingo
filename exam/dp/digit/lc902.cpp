#include <bits/stdc++.h>
using namespace std;


const int N = 35;
const int M = 3;
unordered_map<int, bool> mp;
string nums;
int dp[N][M][M];
int f(int i, int mask, bool isLimit, bool isLeadZero);
void printSpace(int i) {
    for (int x=0; x<i; x++) {
        printf(" ");
    }
}

/**
给定一个按非递减顺序排列的数字数组digits。你可以用任意次数digits[i]来写的数字。

例如，如果digits = ['1','3','5']，我们可以写数字，如'13','551', 和'1351315'。

返回可以生成的小于或等于给定整数 n 的正整数的个数。
*/
int atMostNGivenDigitSet(vector<string>& digits, int n) {
    nums = to_string(n);
    memset(dp, -1, sizeof(dp));
    mp.clear();
    for (auto s : digits) {
        mp[int(s[0] - '0')] = true;
    }
    return f(0, 0, true, true);
}

/*
f(i, mask, isLimit, isLeadZero):
    i: 位置i
    mask: 上一位是啥
    isLimit: 之前所有访问的前(i-1)位都命中了上界。
        eg: n=124, 如果前面访问了12, 第3位可选的就只有 [0,4] 了
    isLeadZero: 之前所有访问的前(i-1)位都命中了下界。
        eg: n=124, 如果前面访问了00, 第3位0,1,2,3,4 只用统计一次次数
*/
int f(int i, int mask, bool isLimit, bool isLeadZero) {
    if (i == nums.size()) {
        if (isLeadZero) {
            return 0;
        }
        return 1;
    }
    if (dp[i][isLimit][isLeadZero] != -1) {
//    #ifdef debug
//        printSpace(i);
//        printf("$$i:%d isLimit:%d isLeadZero:%d dp[i]:%d\n", i, isLimit, isLeadZero, dp[i]);
//    #endif
        return dp[i][isLimit][isLeadZero];
    }
    #ifdef debug
        printSpace(i);
        printf("=>");
        printf("i:%d isLimit:%d isLeadZero:%d\n", i, int(isLimit), int(isLeadZero));
    #endif
    int ans = 0;
    if (isLeadZero) {
        ans += f(i+1, mask, false, isLeadZero);
    }
    #ifdef debug
        printSpace(i);
        printf("<==");
        printf("i:%d isLeadZero:%d ans: %d\n", i, int(isLeadZero), ans);
    #endif

    int up, down;

    up = isLimit ? nums[i] - '0' : 9;
    down = isLeadZero ? 1 : 0;
    for (int v = down; v <= up; v++) {
        if (!mp[v]) continue;
        int nextMask = mask;
        int nextLimit = (isLimit) && (v == up);
        int tmp = f(i+1, nextMask, nextLimit, false);
        ans += tmp;
        #ifdef debug
        printSpace(i);
        printf("**");
        printf("i:%d v: %d tmp:%d\n", i, v, tmp);
        #endif
    }
    #ifdef debug
        printSpace(i);
        printf("<=");
        printf("i:%d up:%d down:%d ans: %d\n", i, up, down, ans);
    #endif

    dp[i][isLimit][isLeadZero] = ans;
    return ans;
}

int main() {
    cout << "LC 902" << endl;
    auto inputs = vector<vector<string>>{{"1","3","5","7"}, {"7"}, {"1","4","9"}, {"1","4","9"}, {"1","4","9"},{"1","4","9"}};
    auto intputs2 = vector<int>{100, 8, 200, 149, 999, 1000000000};
    vector<int> outputs = vector<int>{20, 1, 21, 18,39, 29523};
    bool allPass = true;

    for (int i=0; i<inputs.size(); i++) {
        auto input = inputs[i];
        auto input2 = intputs2[i];
        int expect = outputs[i];
        int ans = atMostNGivenDigitSet(input, input2);
        if (ans != expect) {
            cout << "case[" << i+1 << "]";
            cout << " input: ";
            for (auto s : input) {
                cout << s << " ";
            }
            cout << endl;
            printf("input2: %d output: %d, but expect: %d\n", input2, ans, expect);
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