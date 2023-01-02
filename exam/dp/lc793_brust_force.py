import functools

def solve(n):
    nums = [i+1 for i in range(n)]
    res = functools.reduce(lambda x,y: x*y, nums)
    s = str(res)
    ll = len(s)
    cnt = 0
    for i in range(ll-1, -1, -1):
        if s[i] == '0':
            cnt += 1
        else:
            break
    return cnt

vis = {}
n = 1000
for i in range(n):
    vis[solve(i+1)] = True

print(vis.keys())