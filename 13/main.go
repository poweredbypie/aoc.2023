package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pattern struct {
	lines []string
}

func (p *Pattern) ReflectValue(dist int) int {
	rowDist := func(row1, row2 int) int {
		dist := 0
		for col := 0; col < len(p.lines[row1]); col += 1 {
			if p.lines[row1][col] != p.lines[row2][col] {
				dist += 1
			}
		}
		return dist
	}
	reflectRowDist := func(row int) int {
		// Check in both directions
		dist := 0
		for iter := 0; row+iter+1 < len(p.lines) && row-iter >= 0; iter += 1 {
			dist += rowDist(row+iter+1, row-iter)
		}
		return dist
	}
	colDist := func(col1, col2 int) int {
		dist := 0
		for row := 0; row < len(p.lines); row += 1 {
			if p.lines[row][col1] != p.lines[row][col2] {
				dist += 1
			}
		}
		return dist
	}
	reflectColDist := func(col int) int {
		dist := 0
		for iter := 0; col+iter+1 < len(p.lines[0]) && col-iter >= 0; iter += 1 {
			dist += colDist(col+iter+1, col-iter)
		}
		return dist
	}
	for i := 0; i < len(p.lines)-1; i += 1 {
		if reflectRowDist(i) == dist {
			// log.Printf("Row %v has a distance of 1 from equivalent", i)
			return (i + 1) * 100
		}
	}
	for i := 0; i < len(p.lines[0])-1; i += 1 {
		if reflectColDist(i) == dist {
			// log.Printf("Col %v has a distance of 1 from equivalent", i)
			return i + 1
		}
	}
	panic("No reflection found")
}

func NewPattern(scan *bufio.Scanner) *Pattern {
	lines := []string{}
	for scan.Scan() {
		line := scan.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return nil
	}
	return &Pattern{lines}
}

func main() {
	file, _ := os.Open("input")
	scan := bufio.NewScanner(file)
	sumPerf := 0
	sumSmudge := 0
	for {
		pattern := NewPattern(scan)
		if pattern == nil {
			break
		}
		sumPerf += pattern.ReflectValue(0)
		sumSmudge += pattern.ReflectValue(1)
	}
	fmt.Printf("Sum of reflect values is %v\n", sumPerf)
	fmt.Printf("Sum of reflect values with 1 smudge is %v\n", sumSmudge)
}
