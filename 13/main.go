package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Pattern struct {
	lines []string
}

func (p *Pattern) ReflectValue() int {
	reflectsAtRow := func(row int) bool {
		// Check in both directions
		for iter := 0; row+iter+1 < len(p.lines) && row-iter >= 0; iter += 1 {
			if p.lines[row+iter+1] != p.lines[row-iter] {
				return false
			}
		}
		return true
	}
	colEquals := func(col1, col2 int) bool {
		for row := 0; row < len(p.lines); row += 1 {
			if p.lines[row][col1] != p.lines[row][col2] {
				return false
			}
		}
		return true
	}
	reflectsAtCol := func(col int) bool {
		for iter := 0; col+iter+1 < len(p.lines[0]) && col-iter >= 0; iter += 1 {
			if !colEquals(col+iter+1, col-iter) {
				return false
			}
		}
		return true
	}
	for i := 0; i < len(p.lines)-1; i += 1 {
		if reflectsAtRow(i) {
			log.Printf("Reflects at row %v", i)
			return (i + 1) * 100
		}
	}
	for i := 0; i < len(p.lines[0])-1; i += 1 {
		if reflectsAtCol(i) {
			log.Printf("Reflects at column %v", i)
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
	sum := 0
	for {
		pattern := NewPattern(scan)
		if pattern == nil {
			break
		}
		val := pattern.ReflectValue()
		log.Printf("Val for pattern %v is %v", pattern, val)
		sum += val
	}
	fmt.Printf("Sum of reflect values is %v\n", sum)
}
