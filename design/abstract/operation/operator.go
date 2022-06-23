package operation

import (
	"context"
	"fmt"
)

var _ IOperator = &BaseOperator{}

type OperateType int

const (
	OperateType_Insert OperateType = iota + 1 // 新增
	OperateType_Update                        // 变更
	OperateType_Delete                        // 删除
)

var _ IOperator = &BaseOperator{}

// IOperator 编辑操作抽象
type IOperator interface {
	// CheckParams 参数校验
	CheckParams() error
	// Check 状态校验
	Check() error
	// GetCurrent 获取草稿版本数据
	GetCurrent() error
	// Diff 提审数据对比草稿数据
	Diff() map[int64]OperateType
	// Save 落库
	Save() error
	// SubmitAudit 提审
	SubmitAudit() error
	// SendEvent 发布变更事件
	SendEvent() error
}

type ISubmitService interface {
	CheckParam() error

	GetCurrent() error

	Diff() map[int64]OperateType

	Save() error
}

type BaseOperator struct {
	ctx context.Context

	services []ISubmitService

	changeMap map[int64]OperateType
}

func NewBaseOperator(ctx context.Context, services []ISubmitService) *BaseOperator {
	return &BaseOperator{
		ctx:      ctx,
		services: services,

		changeMap: make(map[int64]OperateType),
	}
}

func (b *BaseOperator) Operate() (err error) {
	if err = b.CheckParams(); err != nil {
		return err
	}
	if err = b.Check(); err != nil {
		return err
	}

	if err = b.GetCurrent(); err != nil {
		return err
	}
	b.changeMap = b.Diff()

	if err = b.Save(); err != nil {
		return err
	}
	if err = b.SubmitAudit(); err != nil {
		return err
	}
	if err = b.SendEvent(); err != nil {
		return err
	}
	return nil
}

func (b *BaseOperator) CheckParams() (err error) {
	fmt.Println("base CheckParam")
	for _, svc := range b.services {
		if err = svc.CheckParam(); err != nil {
			break
		}
	}
	return
}

func (b *BaseOperator) Check() error {
	fmt.Println("check")
	return nil
}

func (b *BaseOperator) GetCurrent() error {
	for _, svc := range b.services {
		if err := svc.GetCurrent(); err != nil {
			return err
		}
	}
	fmt.Println("GetCurrent done")
	return nil
}

func (b *BaseOperator) Diff() map[int64]OperateType {
	changeMap := make(map[int64]OperateType)
	for _, svc := range b.services {
		svcDiffMap := svc.Diff()
		for k, _ := range svcDiffMap {
			changeMap[k] = svcDiffMap[k]
		}
	}
	fmt.Println("Diff")
	return changeMap
}

func (b *BaseOperator) Save() (err error) {
	for k, v := range b.changeMap {
		fmt.Printf("qId %v operateType %v\n", k, v)
	}
	fmt.Println("base Save")
	return
}

func (b *BaseOperator) SubmitAudit() (err error) {
	fmt.Println("base SubmitAudit")
	return
}

func (b *BaseOperator) SendEvent() (err error) {
	fmt.Println("base SendEvent")
	return
}

type SubjectService struct{}

func (s *SubjectService) CheckParam() error {
	fmt.Println("SubjectService CheckParams")
	return nil
}

func (s *SubjectService) GetCurrent() error {
	fmt.Println("SubjectService GetCurrent")
	return nil
}

func (s *SubjectService) Diff() map[int64]OperateType {
	return map[int64]OperateType{1: OperateType_Insert}
}

func (s *SubjectService) Save() error {
	fmt.Println("SubjectService Save")
	return nil
}

type BackgroundService struct{}

func (s *BackgroundService) CheckParam() error {
	fmt.Println("BackgroundService CheckParams")
	return nil
}

func (s *BackgroundService) GetCurrent() error {
	fmt.Println("BackgroundService GetCurrent")
	return nil
}

func (s *BackgroundService) Diff() map[int64]OperateType {
	return map[int64]OperateType{2: OperateType_Update}
}

func (s *BackgroundService) Save() error {
	fmt.Println("BackgroundService Save")
	return nil
}
