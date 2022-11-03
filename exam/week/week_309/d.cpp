typedef pair<int, int> pii;
class Solution {
public:
    int mostBooked(int n, vector<vector<int>>& meetings) {
        auto f = [](bool)(pii x, pii y) {
            if (x.first == y.first) {
                return x.end > y.end;
            }
            return x.first > y.first;
        };
        priority_queue<pii, vector<pii>, decltype(f)> que(f);
        vector<int> cnt(n, 0);
        vector<int> last(n, 0);

        for (auto &meet : meetings) {
            que.push({meet[0], meet[1]});
        }

        int l = meetings.size();

        while (!que.empty()) {
            pii tmp = que.top(); que.pop();
            int pos = 0;
            for (int j=1; j<n; j++) {
                if (last[pos] > last[j]) {
                    pos = j;
                }
            }
            if (tmp.first >= last[pos]) {
                last[pos] = tmp.second;
                cnt[pos]++;
            } else {
                int diff = last[pos] - tmp.first;
                pii tmp2 = {tmp.first + diff, tmp.second + diff};
                que.push(tmp2);
            }

        }


        int pos = 0;
        for (int i=1; i<n; i++) {
            if (cnt[pos] < cnt[i]) {
                pos = i;
            }
        }

        return pos+1;
    }
}