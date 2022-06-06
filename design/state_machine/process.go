package state_machine

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

var (
	debug           bool
	once            sync.Once
	statusToProcess map[StateStatus]IProcessor
	pool            *TaskPool
)

const (
	FailTh        = 10  // 状态机流转过程中，失败比例= FailTh / 100
	ConstTaskPool = 100 // 默认任务池大小
	RetryTimeTh   = 15  // 15次重试失败，表示系统肯定有问题，需要人工介入
)

func Init() {
	once.Do(func() {
		// 注册 status -> processor
		statusToProcess = make(map[StateStatus]IProcessor)
		statusToProcess[StateStatus_Init] = &StoreProcessor{}
		statusToProcess[StateStatus_Assemble] = &AssembleProcessor{}
		statusToProcess[StateStatus_HumanAudit] = &HumanProcessor{}
		statusToProcess[StateStatus_AuditFinished] = &AuditFinishedProcessor{}

		// 初始化异步任务池
		pool = NewTaskPool(ConstTaskPool)
		go pool.Start()
	})
}

type MainProcessor struct {
	node      *Node
	processor IProcessor
}

func NewMainProcessor(node *Node) *MainProcessor {
	obj := &MainProcessor{node: node, processor: nil}
	if p, ok := statusToProcess[node.stateStatus]; !ok {
		panic(fmt.Sprintf("no found %v -> processor", node.stateStatus))
	} else {
		obj.processor = p
	}
	return obj
}

func (m *MainProcessor) Process(ctx context.Context) error {
	if m.processor == nil {
		panic(fmt.Sprintf("taskId %v stateStatus %v no processor", m.node.taskId, m.node.stateStatus))
	}
	//if debug {
	fmt.Printf("taskId %v processor %T stateStatus %v\n", m.node.taskId, m.processor, m.node.stateStatus)
	//}

	if err := m.processor.Check(m.node); err != nil { // 校验不通过的，没法通过重试解决
		fmt.Printf("[AlarmCheckFail] taskId %v stateStatus %v check err: %v\n", m.node.taskId, m.node.stateStatus, err) // 报警
		return err
	}

	nextNode, err := m.processor.Process(ctx, m.node)
	if err != nil {
		return retryProcess(ctx, m.node)
	}
	if nextNode != nil {
		err = NewMainProcessor(nextNode).Process(ctx)
	}
	return err
}

func retryProcess(ctx context.Context, node *Node) error {
	node.retryTimes++
	if node.retryTimes > RetryTimeTh {
		fmt.Printf("[AlarmRetryTooMany] taskId %v stateStatus %v has retry too many\n", node.taskId, node.stateStatus) // 报警
	} else {
		pool.AddNode(ctx, node)
	}
	if debug {
		fmt.Printf("taskId %v stateStatus %v retryTime %d\n", node.taskId, node.stateStatus, node.retryTimes)
	}

	return nil
}

func mockError() error {
	// 模拟随机故障，暂定阈值是 failTh%
	randV := rand.Intn(100)
	if randV < FailTh {
		return errors.New("随机故障")
	}
	return nil
}
