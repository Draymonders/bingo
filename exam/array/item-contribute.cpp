#include <bits/stdc++.h>
using namespace std;
typedef long long ll;

/*
求元素贡献，相关题目
*/

bool debug;

/**
判断密码强度,例如"good",所有子串为
     g o o d go oo od goo ood good
     1 1 1 1 2  1  2  2   2    3
    求和为16，返回16
*/
ll solve(string s) {
    int n = s.size();
    ll ans = 0;
    ll dp[n];   memset(dp, 0, sizeof(dp)); // dp[i]: 位置i，能产生的贡献
    unordered_map<char, int> mp; // key: char val: pos

    for (int i=0; i<n; i++) {
        if (mp.count(s[i])) {
            int pos = mp[s[i]];
            dp[i] += i - pos;
        } else {
            dp[i] = i + 1;
        }
        if (i > 0)
            dp[i] += dp[i-1];

        if (debug) {
            printf("i:%d, dp[i]: %lld\n", i, dp[i]);
        }
        ans += dp[i];
        mp[s[i]] = i;
    }
    return ans;
}

int main() {
    debug = false;
    vector<string> strs = vector<string>{"good", "goods", "abc"};
    vector<ll> expects = vector<ll>{16, 29, 10};

    for (int i=0; i<strs.size(); i++) {
        string s = strs[i];
        ll expected = expects[i];

        ll v = solve(s);
        if (v != expected) {
            cout << "str: " << s << endl;
            printf("v: %lld, but expected: %lld\n", v, expected);
            exit(1);
        }
    }
    cout << "success" << endl;
    return 0;
}
