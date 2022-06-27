package async

import "testing"

func Test_Go(t *testing.T) {
	{
		r := &req{Param: 1}
		asyncGo(r)
	}

	{
		r := &req{Param: 2}
		rsp, err := syncGo(r)

		t.Logf("rsp %v err: %v", rsp, err)
	}
}
