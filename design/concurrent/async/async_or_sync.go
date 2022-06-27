package async

import "time"

type token struct{}

type req struct {
	Param int
}

type resp struct {
	Body int
}

type call struct {
	req  *req
	resp *resp
	err  error
	done chan token
}

// 由上层去判断是否走异步还是串行

// 异步的话 无需关心返回值
func asyncGo(req *req) {
	_go(req, make(chan token, 1))
}

// handleSync 同步方法
func syncGo(req *req) (*resp, error) {
	c := _go(req, make(chan token, 1))
	<-c.done
	return c.resp, c.err
}

func _go(req *req, doneCh chan token) *call {
	rsp := new(resp)
	_call := new(call)
	_call.req = req
	_call.resp = rsp
	_call.done = doneCh

	go func(c *call) {
		time.Sleep(time.Second * 1)
		c.resp.Body = c.req.Param + 1
		c.done <- token{}
	}(_call)
	return _call
}
