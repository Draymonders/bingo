package utils

import "testing"

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
