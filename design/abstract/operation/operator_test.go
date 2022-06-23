package operation

import (
	"context"
	"testing"
)

func Test_Operator(t *testing.T) {
	ctx := context.Background()
	operator := NewBaseOperator(ctx, []ISubmitService{
		&SubjectService{},
		&BackgroundService{},
	})

	if err := operator.Operate(); err != nil {
		t.FailNow()
	}

}
