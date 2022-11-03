package utils

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func Test_Wrap(t *testing.T) {
	fmt.Println(NewWrapper())
}

func TestContain(t *testing.T) {
	OpenDebug()

	{
		i := 1
		arr := []int{1, 2, 3}
		if f := Contain(i, arr); !f {
			t.FailNow()
		}
	}

	{
		i := 1
		mp := map[int]bool{1: true}

		if f := Contain(i, mp); !f {
			t.FailNow()
		}
	}

	{
		i := 1
		arr := []int64{1, 2, 3}
		if f := Contain(i, arr); f {
			t.FailNow()
		}
	}

}

func Test_ToBytes(t *testing.T) {

	obj := SimpleMeta(10, 2, 3)
	t.Logf("%+v", obj)

	buf := (*[valEnd]byte)(unsafe.Pointer(obj))[:valEnd]
	buf[keyEnd] = byte(9)

	t.Logf("%+v", obj)
}

func Test_GetBytes(t *testing.T) {
	obj := &meta{
		HeaderSize: 10,
		KeySize:    2,
		ValSize:    3,
	}
	bytes := (*[valEnd]byte)(unsafe.Pointer(obj))[:valEnd]

	t.Logf("%+v", bytes)
}

func Test_Limiter(t *testing.T) {
	l := NewQpsLimiter(time.Second, 1)

	_ = l.Acquire()
	for i := 1; i < 10; i++ {
		if l.Acquire() {
			t.Logf("err: acquire!!")
			t.FailNow()
		}
	}
	time.Sleep(time.Second)
	if !l.Acquire() {
		t.Logf("err: not acquire!!")
		t.FailNow()
	}
}

func Test_Id_Unmarshal(t *testing.T) {
	str := "{\"poi_id_v2\":6980984069106436127}"
	data := map[string]interface{}{}
	err := JsonToObj(str, &data)
	if err != nil {
		fmt.Println("jsonToObj err")
		return
	}
	fmt.Println(data)
	for _, v := range data {
		fmt.Printf("%v", v)
	}
	//	map[poi_id_v2:6.980984069106436e+18]
	//	6980984069106436096.000000
}

func Test_Uint64_To_Float64(t *testing.T) {
	// 这种情况下会出现不一致
	v := float64(uint64(6980984069106436127))
	fmt.Println(v)
	fmt.Printf("%f\n", v)
	// 6.980984069106436e+18
	// 6980984069106436096.000000
}

func Test_Recover(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			Recover("recover panic test")
			wg.Done()
		}()
		panic("wow~ go func panic")
	}()

	wg.Wait()

}

type Node struct {
	//a int
	B int
	// c *int
	D NodeA
	E *NodeB
}

type NodeA struct {
	V string
}

type NodeB struct {
	v int
}

func buildNode(initValue int, str string) *Node {
	return &Node{
		//a: initValue,
		B: initValue + 1,
		// c: &initValue,
		D: NodeA{
			V: str,
		},
		E: &NodeB{
			v: initValue * 2,
		},
	}
}

func Test_DiffData(t *testing.T) {
	n1 := buildNode(1, "233")
	n2 := buildNode(1, "233")
	var equal bool
	var differentTypeNames []string

	if equal, differentTypeNames = DeepEqual(n1, n2); !equal {
		t.Logf("differentTypeNames %v", differentTypeNames)
		t.FailNow()
	} else {
		t.Logf("equal differentTypeNames %v", differentTypeNames)
	}

	n2.B = 3
	if equal, differentTypeNames = DeepEqual(n1, n2); equal {
		t.Logf("equal differentTypeNames %v", differentTypeNames)
		t.FailNow()
	} else {
		t.Logf("not equal differentTypeNames %v", differentTypeNames)
	}
	n2.B = 2

	n2.E.v = 233
	if equal, differentTypeNames = DeepEqual(n1, n2); equal {
		t.Logf("equal differentTypeNames %v", differentTypeNames)
		t.FailNow()
	} else {
		t.Logf("not equal differentTypeNames %v", differentTypeNames)
	}

}

type INode struct {
	A IFieldA
	B *IFieldA
}

type IFieldA struct {
	V int
}

func buildINode(initVal int) *INode {
	return &INode{A: IFieldA{V: initVal}, B: &IFieldA{V: initVal}}
}

func (x IFieldA) Equal(y interface{}) bool {
	if yStruct, ok := y.(IFieldA); ok {
		return x.V == yStruct.V
	}
	if _, ok := y.(*IFieldA); ok {
		return true
	}
	return true
}

func Test_InterfaceEqual(t *testing.T) {
	x := buildINode(1)
	y := buildINode(2)

	if equal, differentTypeNames := DeepEqual(x, y); !equal {
		t.Logf("not equal differentTypeNames %v", differentTypeNames)
		t.FailNow()
	} else {
		t.Logf("equal differentTypeNames %v", differentTypeNames)
	}
}

func Test_Reflect_Demo1(t *testing.T) {
	m := buildINode(1)

	tt := reflect.TypeOf(m)
	v := reflect.ValueOf(m)

	fmt.Printf("t: %s, kind: %s\nv: %s, kind: %s\n", tt, tt.Kind(), v, v.Kind())
}

func Test_PlaceHolder(t *testing.T) {
	str := "${A} = 1 && ${2B2} = 2 && ${3}"
	vars, err := FactorPlaceholderReplacer.FindPlaceholder(str)
	if err != nil {
		t.FailNow()
	}
	fmt.Println(len(vars), vars)

	newStr, err := FactorPlaceholderReplacer.Replace(str, map[string]interface{}{"A": 1, "2B2": 233})
	if err != nil {
		t.FailNow()
	}
	fmt.Println("newStr", newStr)
}
