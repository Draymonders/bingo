package utils

import "testing"

func Test_FormatFloat64(t *testing.T) {
	{
		bit := 2
		args := []float64{0.012, 0.015, 0.016, 0.4288687420754218, 0.9000589489688305, 0.7275542428109711, -0.03, -0.021, -0.026, 0.0, 0.104, -0.104, -0.106}
		wants := []float64{0.01, 0.02, 0.02, 0.43, 0.90, 0.73, -0.03, -0.02, -0.03, 0.00, 0.10, -0.10, -0.11}

		for i := 0; i < len(args); i++ {
			arg := args[i]
			want := wants[i]
			v := FormatFloat64(arg, bit)
			if v != want {
				t.Logf("arg %v, want %v, but Format64 is %v", arg, want, v)
				t.FailNow()
			}
		}
	}

	{
		bit := 4
		args := []float64{0.01213, 0.01554, 0.01655, 0.8437411942684759, 0.29714084929519585, 0.7201089206044288}
		wants := []float64{0.0121, 0.0155, 0.0166, 0.8437, 0.2971, 0.7201}

		for i := 0; i < len(args); i++ {
			arg := args[i]
			want := wants[i]
			v := FormatFloat64(arg, bit)
			if v != want {
				t.Logf("arg %v, want %v, but Format64 is %v", arg, want, v)
				t.FailNow()
			}
		}
	}
}
