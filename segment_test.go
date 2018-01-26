package overlap

import (
	"testing"
	"fmt"
)

func TestNewLineSegment(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Fatal(r)
		}
	}()

	lines := [][]string{
		{"3", "0", "37.336", "55.7346", "16932", "37.3309", "55.7451", "18976", "WALKING", "3", "1"},
		{"6", "4", "37.3447", "55.7726", "37294", "37.3389", "55.7691", "38492", "WALKING", "0", "0"},
		{"9", "0", "37.3387", "55.7278", "43200", "37.3377", "55.7281", "43638", "DRIVING", "1", "1"},
		{"9", "0", "37.3377", "55.7281", "43638", "37.3384", "55.7324", "44076", "DRIVING", "1", "1"},
		{"9", "0", "37.3384", "55.7324", "44076", "37.3389", "55.7343", "44514", "DRIVING", "1", "1"},
	}

	var segments []*LineSegment
	for idx, l := range lines {
		ls, err := NewLineSegment(idx, l)
		if err != nil {
			t.Fatal(err)
		}
		segments = append(segments, ls)
	}

	for _, v := range segments {
		fmt.Println(v, ">"+v.getLineFormula()+"<")
	}

	// y=1.4x-3.42
	l1 := LineSegment{
		startX: 5.64074,
		startY: 4.47704,
		endX:   0.82952,
		endY:   -2.2586,
	}
	l2 := LineSegment{
		startX: 10.5625,
		startY: 11.3676,
		endX:   -9.1247,
		endY:   -16.1946,
	}
	fmt.Println(l1.getLineFormula(), l2.getLineFormula())
	fmt.Println(l1.getLineFormula() == l2.getLineFormula())

	l3 := LineSegment{
		startX: 5,
		startY: 1,
		endX:   5,
		endY:   13,
	}
	fmt.Println(l3.getLineFormula())

	l4 := LineSegment{
		startX: 3,
		startY: 1,
		endX:   5,
		endY:   1,
	}
	fmt.Println(l4.getLineFormula())
}
