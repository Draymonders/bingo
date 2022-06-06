package concurrent

import "testing"

func TestParallelPrintOddAndEvenV1(t *testing.T) {
	ParallelPrintOddAndEvenV1(10)
	ParallelPrintOddAndEvenV1(9)
}

func TestParallelPrintOddAndEvenV2(t *testing.T) {
	ParallelPrintOddAndEvenV2(100)
	ParallelPrintOddAndEvenV2(99)
}
