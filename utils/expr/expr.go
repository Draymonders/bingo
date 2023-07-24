package expr

import (
	"fmt"
	"reflect"

	"github.com/Knetic/govaluate"
)

type IExpr interface {
	Bool(input map[string]interface{}) (bool, error)
}

type LocalExprImpl struct {
	exprStr string
}

func NewLocalExprImpl(exprStr string) IExpr {
	return &LocalExprImpl{exprStr: exprStr}
}

func (s *LocalExprImpl) Bool(input map[string]interface{}) (bool, error) {
	if len(s.exprStr) == 0 {
		return false, nil
	}

	expr, err := govaluate.NewEvaluableExpression(s.exprStr)
	if err != nil {
		return false, err
	}
	res, err := expr.Evaluate(input)
	if err != nil {
		return false, err
	}
	switch res.(type) {
	case bool:
		return res.(bool), nil
	default:
		fmt.Println("res is not bool")
		return false, nil
	}
}

const IntExprErrVal int64 = -9999

func IntExpr(exprStr string, input map[string]interface{}) (int64, error) {
	if len(exprStr) == 0 {
		return IntExprErrVal, nil
	}

	expr, err := govaluate.NewEvaluableExpression(exprStr)
	if err != nil {
		fmt.Println("err1")
		return IntExprErrVal, err
	}
	res, err := expr.Evaluate(input)
	if err != nil {
		fmt.Println("err2")
		return IntExprErrVal, err
	}
	switch res.(type) {
	case int:
		return int64(res.(int)), nil
	case int64:
		return res.(int64), nil
	default:
		fmt.Println("res is not int, but res.(type)", reflect.TypeOf(res))
		return IntExprErrVal, nil
	}
}
