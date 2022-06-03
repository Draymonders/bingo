package state_machine

import (
	"context"
	"testing"
	"time"
)

// 状态机流转
func Test_StateMachine(t *testing.T) {
	Init()
	ctx := context.Background()

	n := 50
	var nodes []*Node
	for i := 1; i <= n; i++ {
		node := &Node{
			taskId:      int64(10000 + i),
			stateStatus: StateStatus_Init,
			retryTimes:  0,
		}
		nodes = append(nodes, node)
		if err := NewMainProcessor(node).Process(ctx); err != nil {
			t.FailNow()
		}
	}

	// go test -v -count=1 -run Test_StateMachine > out
	// grep -E "stateStatus 400" out | grep -E -o "taskId 1[0-9]+" | cut -d' ' -f 2 | sort | uniq | wc -l
	// 等于 n
	time.Sleep(5 * time.Second)
}
