package utils

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

func init() {
	Init(logrus.DebugLevel)
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
		if IsDebug() {
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

// todo @yubing 增加单测
func Struct2MapJson(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = reflect.TypeOf(v.Interface())
	}

	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		if v.IsNil() {
			return nil
		}

		l := v.Len()
		if l == 0 {
			return nil
		}
		av := make([]interface{}, l)
		for i := 0; i < l; i++ {
			if v.Index(i).IsNil() { // slice 成员，nil 判断
				continue
			}
			r := Struct2MapJson(v.Index(i).Elem().Interface())
			if r != nil {
				av[i] = r
			}
		}
		return av
	} else if v.Kind() == reflect.Map {
		r := make(map[string]interface{}, len(v.MapKeys()))
		for _, element := range v.MapKeys() {
			//fmt.Println(key, element) // how to get the value?
			val := v.MapIndex(element)
			r[element.Interface().(string)] = Struct2MapJson(val.Interface())
		}
		if len(r) == 0 {
			return nil
		}
		return r
	} else if v.Kind() == reflect.Struct {
		var data = make(map[string]interface{}, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			name := t.Field(i).Tag.Get("json")
			if name == "" {
				name = t.Field(i).Name
			}
			v1 := v.Field(i).Interface()
			if v1 == nil {
				continue
			}
			r := Struct2MapJson(v1)
			if r != nil {
				data[name] = r
			}
		}
		return data
	}
	return v.Interface()
}
