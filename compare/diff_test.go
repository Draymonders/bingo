package compare

import (
	"fmt"
	"reflect"
	"testing"
)

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