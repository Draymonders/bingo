typedef long long ll;

/*
6098. 统计得分小于 K 的子数组数目
https://leetcode.cn/problems/count-subarrays-with-score-less-than-k/

范围: 1 <= len(nums) <= 1e5
枚举每个pos, 计算以pos为子数组开始位置的合法的子数组数目
*/
class Solution {
public:
    long long countSubarrays(vector<int>& nums, long long k) {
        int n = nums.size();
        vector<ll> sum(n+1);
        for (int i=1; i<=n; i++)
            sum[i] = sum[i-1] + nums[i-1];
        ll res = 0;
        for (int i=1; i<=n; i++) {
            ll v = calc(sum, n, i, k);
            // printf("i:%d v: %lld\n", i, v);
            res += v;
        }
        return res;
    }


    ll calc(vector<ll>& sum, int n, int pos, ll k) {
        int l = pos, r = n;
        int ans = -1;
        while (l <= r) {
            int m = (l + r) / 2;
            if ((sum[m] - sum[pos-1]) * (m - pos + 1) < k) {
                ans = m, l = m + 1;
            } else {
                r = m-1;
            }
        }
        if (ans == -1 || ans < pos) {
            return 0;
        }
        return ans - pos + 1;
    }
};