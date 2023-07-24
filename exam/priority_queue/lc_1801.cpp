#include <functional>
#include <queue>
#include <vector>
#include <iostream>
using namespace std;


int mod = 1e9 + 7;
struct node {
    int price, amount, type_;
    node(int p, int a, int t) : price(p), amount(a), type_(t) {}
};

bool debug;
// cpp pq 使用

class Solution {
public:
    int getNumberOfBacklogOrders(vector<vector<int>>& orders) {
        auto buyFunc = [](node a, node b) {
            return a.price < b.price;
        };
        auto sellFunc = [](node a, node b) {
            return a.price > b.price;
        };

        priority_queue<node, vector<node>, decltype(buyFunc)> buyQueue(buyFunc);
        priority_queue<node, vector<node>, decltype(sellFunc)> sellQueue(sellFunc);


       int i = 0;

        for (auto order : orders) {
            int curPrice = order[0], curAmount = order[1], curType = order[2];
            if (curType == 0) { // buy
                while (!sellQueue.empty() && curAmount > 0) {
                    if (sellQueue.top().price <= curPrice) {
                        node tmp = sellQueue.top();
                        sellQueue.pop();

                        int mn = min(curAmount, tmp.amount);
                        tmp.amount -= mn;
                        curAmount -= mn;
                        if (debug)
                          printf("[buy] curPrice %d may use temp price %d mn %d\n", curPrice, tmp.price, mn);
                        if (tmp.amount > 0) {
                            sellQueue.push(tmp);
                        }
                    } else {
                        break;
                    }
                }
                if (curAmount > 0) {
                    buyQueue.push(node(curPrice, curAmount, curType));
                }
            } else {
               while (!buyQueue.empty() && curAmount > 0) {
                   if (buyQueue.top().price >= curPrice) {
                       node tmp = buyQueue.top();
                       buyQueue.pop();

                       int mn = min(curAmount, tmp.amount);
                       tmp.amount -= mn;
                       curAmount -= mn;

                       if (debug)
                         printf("[sell] curPrice %d may use temp price %d mn %d\n", curPrice, tmp.price, mn);
                       if (tmp.amount > 0) {
                           buyQueue.push(tmp);
                       }
                   } else {
                       break;
                   }
               }
               if (curAmount > 0) {
                   sellQueue.push(node(curPrice, curAmount, curType));
               }
            }

            if (debug) {
            printf("i: %d\n", i++);
               vector<node> buyTemps, sellTemps;
               while (!buyQueue.empty()) {
                   buyTemps.push_back(buyQueue.top());
                   buyQueue.pop();
               }
               while (!sellQueue.empty()) {
                   sellTemps.push_back(sellQueue.top());
                   sellQueue.pop();
               }
               printf("=== buy queue == \n");
               for (auto b : buyTemps) {
                   printf("price %d amount %d\n", b.price, b.amount);
               }
               printf("=== sell queue===\n");
               for (auto s : sellTemps) {
                   printf("price %d amount %d\n", s.price, s.amount);
               }

               for (auto b : buyTemps) {
                   buyQueue.push(b);
               }
               for (auto s : sellTemps) {
                   sellQueue.push(s);
               }
               printf("\n");
            }
        }

        int res = 0;
        while (!buyQueue.empty()) {
            node tmp = buyQueue.top();
            res = (res + tmp.amount) % mod ;
            buyQueue.pop();
        }
        while (!sellQueue.empty()) {
            node tmp = sellQueue.top();
            res = (res + tmp.amount) % mod;
            sellQueue.pop();
        }
        return res;
    }
};



int main()
{

    debug = false;
    Solution s;

    vector< vector<vector<int>> > ins = {
        {{10,5,0},{15,2,1},{25,1,1},{30,4,0}},
        {{7,1000000000,1},{15,3,0},{5,999999995,0},{5,1,1}}
    };

    vector<int> outs = {
        6,
        999999984
    };
    int i = 0;
    for (auto _in : ins) {
        int got = s.getNumberOfBacklogOrders(_in);
        if (got != outs[i]) {
            printf("want %d, but got %d\n", outs[i], got);
            exit(1);
        }
        i++;
    }
    return 0;
    // todo @yubing test [[10,5,0],[15,2,1],[25,1,1],[30,4,0]] => {xxx}
}