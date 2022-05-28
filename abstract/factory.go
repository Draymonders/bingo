package abstract

import (
	. "github.com/draymonders/bingo/log"
)

/**
工厂方法实现，根据不同的操作类型，实现对应的判断逻辑
*/

type OperateType int64

const (
	OperateType_Settle         OperateType = iota // 入驻提审
	OperateType_ModifyCategory                    // 修改类目
)

type IQuickAudit interface {
	IsHit() bool     // 是否命中
	GetName() string // 方法名
	EmitMetric()     // 程序打点
}

func GetQuickAuditByOperateType(operateType OperateType) IQuickAudit {
	switch operateType {
	case OperateType_Settle:
		return &QuickSettleService{}
	case OperateType_ModifyCategory:
		return &CategoryExemptService{}
	}
	return nil
}

type CategoryExemptService struct{}

func (s *CategoryExemptService) IsHit() bool {
	Log.Info("CategoryExemptService is hit")
	return true
}

func (s *CategoryExemptService) GetName() string {
	return "类目免审"
}

func (s *CategoryExemptService) EmitMetric() {
}

type QuickSettleService struct{}

func (s *QuickSettleService) IsHit() bool {
	Log.Info("QuickSettleService is hit")
	return true
}

func (s *QuickSettleService) GetName() string {
	return "极速入驻"
}

func (s *QuickSettleService) EmitMetric() {
}
