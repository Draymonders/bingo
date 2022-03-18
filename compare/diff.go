package compare

/*
   1. 按字段一个一个的去diff，后续维护比较麻烦
   2. 某些字段支持特殊的逻辑 （对应实现ICompare接口）
   3. 确保传入的参数类型一致，同为指针 or 同为结构体
   4. 一级参数必须都为 public的，即字段参数名必须大写开头，否则调用会panic！
*/

import (
	"reflect"

	. "github.com/draymonders/bingo/log"
	"github.com/sirupsen/logrus"
)

var debug bool

func init() {
	Init(logrus.DebugLevel)
	debug = false
}

func SetDebug(f bool) {
	debug = true
}

type ICompare interface {
	Equal(interface{}) bool
}

// 建议特殊字段实现上面这个ICompare接口
func DeepEqual(x, y interface{}) (bool, []string) {
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		// 类型都不一样
		Log.Warnf("v1Type=%+v, v2Type=%+v, x=%+v, y=%+v", v1.Type(), v2.Type(), x, y)
		return false, nil
	}

	t := reflect.TypeOf(x)
	switch v1.Kind() {
	case reflect.Ptr:
		return deepEqualWithField(v1.Elem(), v2.Elem(), t.Elem())
	case reflect.Struct:
		return deepEqualWithField(v1, v2, t)
	default:
		// 不支持
		Log.Warnf("kind=%+v, x=%+v, y=%+v", v1.Kind(), x, y)
		return false, nil
	}
}

func deepEqualWithField(v1, v2 reflect.Value, t reflect.Type) (bool, []string) {
	diffNames := make([]string, 0)

	// 遍历所有的字段
	for i := 0; i < v1.NumField(); i++ {
		// field为type StructField类型
		field := v1.Field(i)
		if debug {
			Log.Debugf("t.Field(%v).Name %v v1.Field().Interface() %v v2.Field().Interface() %v ", i, t.Field(i).Name, field.Interface(), v2.Field(i).Interface())
		}


		if !isEqualInterface(field.Interface(), v2.Field(i).Interface()) {
			diffNames = append(diffNames, t.Field(i).Name)
		}
	}

	return len(diffNames) == 0, diffNames
}

func isEqualInterface(x, y interface{}) bool {
	if x, ok := x.(ICompare); ok {
		if x.Equal(y) {
			return true
		}
		return false
	}
	return reflect.DeepEqual(x, y)
}
