package search

import "testing"

func Test_canFinish(t *testing.T) {
	//{
	//	numCourses := 2
	//	prerequisites := [][]int{{1, 0}}
	//	want := true
	//	got := canFinish(numCourses, prerequisites)
	//	if got != want {
	//		t.Logf("got %v want %v", got, want)
	//		t.FailNow()
	//	}
	//}
	//{
	//	numCourses := 2
	//	prerequisites := [][]int{{1, 0}, {0, 1}}
	//	want := false
	//	got := canFinish(numCourses, prerequisites)
	//	if got != want {
	//		t.Logf("got %v want %v", got, want)
	//		t.FailNow()
	//	}
	//}
	//{
	//	numCourses := 4
	//	prerequisites := [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 1}}
	//	want := false
	//	got := canFinish(numCourses, prerequisites)
	//	if got != want {
	//		t.Logf("got %v want %v", got, want)
	//		t.FailNow()
	//	}
	//}

	{
		// bad case
		numCourses := 20
		prerequisites := [][]int{{0, 10}, {3, 18}, {5, 5}, {6, 11}, {11, 14}, {13, 1}, {15, 1}, {17, 4}}
		want := false
		got := canFinish(numCourses, prerequisites)
		if got != want {
			t.Logf("got %v want %v", got, want)
			t.FailNow()
		}
	}

}
