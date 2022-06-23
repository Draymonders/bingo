package stock

import "testing"

func Test_Init(t *testing.T) {
	Init()
}

func Test_V1(t *testing.T) {
	IncrementV1(1)
}

func Test_V2(t *testing.T) {
	IncrementV2(1)
}
