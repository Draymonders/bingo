package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	expr := "${a} == 1 && ${b} == 2 && ${c} < 3"
	// 哨兵上绑定的 因子
	factorNames := []string{"a", "b", "c"}
	factorVals := []interface{}{1, 2, 2}

	for i, name := range factorNames {
		varName := fmt.Sprintf("${%s}", name)
		factorVal := factorVals[i]

		expr = strings.ReplaceAll(expr, varName, toString(factorVal))
	}
	fmt.Println("=====")

	exprAst, err := parser.ParseExpr(expr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fset := token.NewFileSet()
	ast.Print(fset, exprAst)
	fmt.Println("=====")

	flag := Parse(exprAst)
	if !flag {
		os.Exit(1)
	}
}

func Parse(exprAst ast.Expr) bool {
	if exprAst == nil {
		panic("exprAst is nil")
	}
	switch exprAst.(type) {
	case *ast.BinaryExpr:
		_ast := exprAst.(*ast.BinaryExpr)
		return dfs(_ast)
	default:
		panic(fmt.Sprintf("not support exprAst.Type %v", reflect.TypeOf(exprAst)))
	}

}

// dfs下 判定是否是符合条件
func dfs(_ast *ast.BinaryExpr) bool {
	switch _ast.Op {
	case token.LAND: // && 操作
		v1 := Parse(_ast.X)
		if !v1 {
			return false
		}
		v2 := Parse(_ast.Y)
		if !v2 {
			return false
		}
		return true
	case token.EQL: // == 操作
		// 左右值解析
		v1 := _ast.X.(*ast.BasicLit)
		v2 := _ast.Y.(*ast.BasicLit)
		if v1.Kind == v2.Kind && v1.Value == v2.Value {
			return true
		}
		return false
	case token.LSS:
		v1 := _ast.X.(*ast.BasicLit)
		v2 := _ast.Y.(*ast.BasicLit)
		if v1.Kind == v2.Kind && v1.Value < v2.Value {
			return true
		}
		return false

	default:
		panic(fmt.Sprintf("not support _ast.Op %v _ast.OpPos %v", _ast.Op, _ast.OpPos))
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case int:
		return strconv.FormatInt(int64(v.(int)), 10)
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	default:
		panic(fmt.Sprintf("v type %v", reflect.TypeOf(v)))
	}
}
