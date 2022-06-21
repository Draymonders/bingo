package factory

import (
	"testing"
)

func Test_QuickAudit(t *testing.T) {
	if service := GetQuickAuditByOperateType(OperateType_Settle); service != nil {
		t.Logf("settle svc %T, name: %v", service, service.GetName())
		service.IsHit()
	}

	if service := GetQuickAuditByOperateType(OperateType_ModifyCategory); service != nil {
		t.Logf("modifyCategory svc %T, name: %v", service, service.GetName())
		service.IsHit()
	}
}
