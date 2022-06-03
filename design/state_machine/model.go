package state_machine

import (
	"context"
	"fmt"
)

const (
	ErrNotSupport = "processor %T only support taskStatus %v, but find %v"
)

// AuditStatus 审核状态
type AuditStatus int

/*
AuditStatus_Auditing ----> AuditStatus_Pass
                      \
                       -> AuditStatus_Reject
*/
const (
	AuditStatus_Auditing AuditStatus = iota // 审核中
	AuditStatus_Pass                        // 通过
	AuditStatus_Reject                      // 驳回
)

// TaskStatus 任务状态
type TaskStatus int

/*
TaskStatus_Init ---> TaskStatus_To_Audit ---> TaskStatus_Pass -------> TaskStatus_CallBack------> TaskStatus_End
                                            \                    /            \
                                             ->  TaskStatus_Reject             -> TaskStatus_DirtyData
*/
const (
	TaskStatus_Init TaskStatus = 100 // 初始化
	//TaskStatus_MachinePass    TaskStatus = 210 // 机审通过
	//TaskStatus_MachineReject  TaskStatus = 220 // 机审通过
	//TaskStatus_NeedHumanAudit TaskStatus = 230 // 机审无法判定，需要人工判定

	TaskStatus_To_Audit  TaskStatus = 300  // 送审成功，审核中
	TaskStatus_Pass      TaskStatus = 310  // 人审通过
	TaskStatus_Reject    TaskStatus = 310  // 人审驳回
	TaskStatus_CallBack  TaskStatus = 800  // 回调外部
	TaskStatus_DirtyData TaskStatus = 910  // 卡审
	TaskStatus_End       TaskStatus = 1000 // 回调完成
)

// StateStatus 状态机的状态
type StateStatus int

/*
StateStatus_Store -> StateStatus_Assemble -> StateStatus_HumanAudit
*/
const (
	StateStatus_Init          StateStatus = 100 // 落库
	StateStatus_Assemble      StateStatus = 200 // 组装送审数据
	StateStatus_HumanAudit    StateStatus = 300 // 人审
	StateStatus_AuditFinished StateStatus = 400 // 人审完成
)

type Node struct {
	taskId      int64       // 提审任务Id
	stateStatus StateStatus // 状态机状态
	retryTimes  int64       // 重试次数
}

type IProcessor interface {
	Check(*Node) error                             // 每个processor 校验下是否状态满足
	Process(context.Context, *Node) (*Node, error) // 处理完当前节点后，下一个节点要处理的节点
}

// StoreProcessor 落库 处理者
type StoreProcessor struct{}

func (s *StoreProcessor) Check(node *Node) error {
	return nil
}

func (s *StoreProcessor) Process(ctx context.Context, node *Node) (*Node, error) {
	return &Node{
		taskId:      node.taskId,
		stateStatus: StateStatus_Assemble,
	}, mockError()
}

// AssembleProcessor 数据组装 处理者
type AssembleProcessor struct{}

func (s *AssembleProcessor) Check(node *Node) error {
	return nil
}

func (s *AssembleProcessor) Process(ctx context.Context, node *Node) (*Node, error) {
	return &Node{
		taskId:      node.taskId,
		stateStatus: StateStatus_HumanAudit,
	}, mockError()
}

// HumanProcessor 人审 处理者
type HumanProcessor struct{}

func (s *HumanProcessor) Check(node *Node) error {
	return nil
}

func (s *HumanProcessor) Process(ctx context.Context, node *Node) (*Node, error) {
	return &Node{
		taskId:      node.taskId,
		stateStatus: StateStatus_AuditFinished,
	}, mockError()
}

// AuditFinishedProcessor 审核完成 处理者
type AuditFinishedProcessor struct{}

func (s *AuditFinishedProcessor) Check(node *Node) error {
	return nil
}

func (s *AuditFinishedProcessor) Process(ctx context.Context, node *Node) (*Node, error) {
	// 终结态了，不需要后续流转了
	return nil, mockError()
}

func checkState(status StateStatus, s IProcessor, shouldStatus StateStatus) error {
	if status != shouldStatus {
		return fmt.Errorf(ErrNotSupport, s, status, shouldStatus)
	}
	return nil
}
