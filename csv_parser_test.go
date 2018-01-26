package overlap

import (
	"testing"
	"fmt"
)

func TestParseCSVFile(t *testing.T) {
	i, err := HandleCSV("segments_1000.csv", true)
	if err != nil {
		t.Fatal(err)
	}
	i2, _ := HandleCSV("segments_1000.csv", false)

	fmt.Printf("Кол-во пар отрезков с наложениями:"+
		" %d - с учётом наложений на свой же путь, %d - и без\r\n", len(i), len(i2))
}
