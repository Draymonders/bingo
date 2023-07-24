package expr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BoolExpr(t *testing.T) {

	a := assert.New(t)
	exprStr := "score >= 0.5"

	res, err := NewLocalExprImpl(exprStr).Bool(map[string]interface{}{"score": 0.666})
	a.Equal(err, nil)
	a.Equal(res, true)

	res, err = NewLocalExprImpl(exprStr).Bool(map[string]interface{}{"score": 0})
	a.Equal(err, nil)
	a.Equal(res, false)
}

func Test_IntExpr(t *testing.T) {
	a := assert.New(t)
	exprStr := "score + 1"

	res, err := IntExpr(exprStr, map[string]interface{}{"score": 2})
	a.Equal(err, nil)
	a.Equal(res, 3)

	res, err = IntExpr(exprStr, map[string]interface{}{"score": -1})
	a.Equal(err, nil)
	a.Equal(res, 0)
}

func Test_BoolVersionExpr(t *testing.T) {
	a := assert.New(t)

	exprStr := "(uid % 100 >= grayRange && v1 == '1.0' && v1_score >= 0.6) || " +
		"(uid % 100 < grayRange && v2 == '2.0' && v2_score < 0.5)"

	{
		res, err := NewLocalExprImpl(exprStr).Bool(map[string]interface{}{
			"uid":       123,
			"grayRange": 0,
			"v1":        "1.0",
			"v1_score":  0.6,
			"v2":        "2.0",
			"v2_score":  0.4,
		})
		a.Equal(err, nil)
		a.Equal(res, true)
	}
	{
		res, err := NewLocalExprImpl(exprStr).Bool(map[string]interface{}{
			"uid":       123,
			"grayRange": 100,
			"v1":        "1.0",
			"v1_score":  0.6,
			"v2":        "2.0",
			"v2_score":  0.4,
		})
		a.Equal(err, nil)
		a.Equal(res, true)
	}

}
