package overlap

import (
	"strconv"
	"math"
	"fmt"
)

type LineSegment struct {
	internalID int
	routeID    int
	startX     float64
	startY     float64
	endX       float64
	endY       float64
	length     float64
	rawData    *[]string
}

const float64BitSize = 64

// решил для упрощения (и в силу малых расстояний) проводить вычисления в евклидовой геометрии
func NewLineSegment(id int, rawData []string) (*LineSegment, error) {
	routeID, err := strconv.Atoi(rawData[0])
	if err != nil {
		return nil, err
	}

	startX, err := strconv.ParseFloat(rawData[2], float64BitSize)
	if err != nil {
		return nil, err
	}
	startY, err := strconv.ParseFloat(rawData[3], float64BitSize)
	if err != nil {
		return nil, err
	}
	endX, err := strconv.ParseFloat(rawData[5], float64BitSize)
	if err != nil {
		return nil, err
	}
	endY, err := strconv.ParseFloat(rawData[6], float64BitSize)
	if err != nil {
		return nil, err
	}

	return &LineSegment{
		internalID: id,
		routeID:    routeID,
		startX:     startX,
		startY:     startY,
		endX:       endX,
		endY:       endY,
		length:     math.Sqrt(math.Pow(startX-endX, 2) + math.Pow(startY-endY, 2)),
		rawData:    &rawData,
	}, nil
}

func almostZero(val1 float64) bool {
	return math.Abs(val1) < 0.00001
}

// точность в три знака после запятой -- ничем не обоснована
func (l *LineSegment) getLineFormula() string {
	deltaX := l.startX - l.endX
	deltaY := l.startY - l.endY

	// Прямая параллельна абциссе
	if almostZero(deltaX) {
		return fmt.Sprintf("y=%.3f", l.startX)
	}

	// Прямая параллельна ординате
	if almostZero(deltaY) {
		return fmt.Sprintf("x=%3.f", l.startY)
	}
	k := deltaY / deltaX
	return fmt.Sprintf("y=%.3f*x +(%.3f)", k, l.startY-l.startX*k)
}

// проверяется нахождения точек первого отрезка внутри или на границах второго отрезка (и точек второго в первом)
func (l *LineSegment) isOverlapped(l2 *LineSegment) bool {
	if l.isInLine(l2.startX, l2.startY) || l.isInLine(l2.endX, l2.endY) {
		return true
	}

	if l2.isInLine(l.startX, l.startY) || l2.isInLine(l.endX, l.endY) {
		return true
	}
	return false
}

func (l *LineSegment) isInLine(x, y float64) bool {
	if x >= math.Min(l.startX, l.endX) && x <= math.Max(l.startX, l.endX) {
		return true
	}
	if y >= math.Min(l.startY, l.endY) && y <= math.Max(l.startY, l.endY) {
		return true
	}
	return false
}
