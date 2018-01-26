package overlap

import (
	"os"
	"encoding/csv"
	"fmt"
	"runtime/debug"
	"encoding/json"
)

func HandleCSV(fp string, checkSelfOverlapping bool) (totalOverlaps Overlaps, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%#v %s", r, debug.Stack())
		}
	}()

	var lines [][]string
	lines, err = ParseData(fp)
	if err != nil {
		return
	}

	m := make(map[string][]*LineSegment)
	var s *LineSegment
	// пропускаем первую строку с описанием
	for idx, line := range lines[1:] {
		s, err = NewLineSegment(idx, line)
		if err != nil {
			return
		}

		// вычисляем формулу прямой, на которой лежит этот отрезок
		// потом ищем наложения у отрезков, лежащих на одной прямой
		k := s.getLineFormula()
		m[k] = append(m[k], s)
	}

	for _, v := range m {
		if len(v) == 1 {
			continue
		}

		overlaps := findOverlaps(v, checkSelfOverlapping)
		if len(overlaps) > 0 {
			totalOverlaps = append(totalOverlaps, overlaps...)
		}
	}
	return
}

func ParseData(fp string) ([][]string, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	return lines, err
}

type Overlaps [][2]*LineSegment

func (o *Overlaps) dump() ([]byte, error) {
	data := make([][2]*[]string, 0, len(*o))
	for _, v := range *o {
		data = append(data, [2]*[]string{v[0].rawData, v[1].rawData})
	}
	b, err := json.MarshalIndent(data, "", "    ")
	return b, err
}

func findOverlaps(segments []*LineSegment, checkSelfOverlapping bool) Overlaps {
	var overlaps Overlaps
	for idx, v1 := range segments[:len(segments)-1] {
		for _, v2 := range segments[idx:] {
			if v1.internalID == v2.internalID || v1.routeID == v2.routeID && !checkSelfOverlapping {
				continue
			}

			if v1.length != v2.length {
				if v1.isOverlapped(v2) {
					overlaps = append(overlaps, [2]*LineSegment{v1, v2})
				}
			}
		}
	}
	return overlaps
}
