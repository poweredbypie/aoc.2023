package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Rocks struct {
	slots [][]byte
}

func NewRocks(file *os.File) *Rocks {
	scan := bufio.NewScanner(file)
	slots := [][]byte{}
	for scan.Scan() {
		slots = append(slots, []byte(scan.Text()))
	}
	return &Rocks{slots}
}

func inc(val int) int {
	return val + 1
}

func dec(val int) int {
	return val - 1
}

func atLeast(min int) func(int) bool {
	return func(val int) bool {
		return val >= min
	}
}

func atMost(max int) func(int) bool {
	return func(val int) bool {
		return val <= max
	}
}

// This sucks but it's kind of generic ish
func (r *Rocks) tiltRock(row, col int, cond func(int) bool, nextRow, nextCol func(int) int) {
	// Start assuming rows
	nextIter := nextRow
	start := row
	charAt := func(iter int) *byte {
		return &r.slots[iter][col]
	}
	// Otherwise assume columns
	if nextRow == nil {
		nextIter = nextCol
		start = col
		charAt = func(iter int) *byte {
			return &r.slots[row][iter]
		}
	}
	// Actual loop
	for iter := start; cond(iter); iter = nextIter(iter) {
		rock := charAt(iter)
		next := charAt(nextIter(iter))
		if *next != '.' {
			return
		}
		*next, *rock = *rock, *next
	}
}

func (r *Rocks) TiltNorth() {
	// Starts at 1 since row 0 is already "tilted"
	for row := 1; row < len(r.slots); row += 1 {
		for col := 0; col < len(r.slots[row]); col += 1 {
			if r.slots[row][col] == 'O' {
				r.tiltRock(row, col, atLeast(1), dec, nil)
			}
		}
	}
}

func (r *Rocks) TiltWest() {
	for col := 1; col < len(r.slots[0]); col += 1 {
		for row := 0; row < len(r.slots); row += 1 {
			if r.slots[row][col] == 'O' {
				r.tiltRock(row, col, atLeast(1), nil, dec)
			}
		}
	}
	return
}

func (r *Rocks) TiltSouth() {
	for row := len(r.slots) - 2; row >= 0; row -= 1 {
		for col := 0; col < len(r.slots[row]); col += 1 {
			if r.slots[row][col] == 'O' {
				r.tiltRock(row, col, atMost(len(r.slots)-2), inc, nil)
			}
		}
	}
	return
}

func (r *Rocks) TiltEast() {
	for col := len(r.slots[0]) - 2; col >= 0; col -= 1 {
		for row := 0; row < len(r.slots); row += 1 {
			if r.slots[row][col] == 'O' {
				r.tiltRock(row, col, atMost(len(r.slots[row])-2), nil, inc)
			}
		}
	}
	return
}

func (r *Rocks) TiltCycle() {
	r.TiltNorth()
	r.TiltWest()
	r.TiltSouth()
	r.TiltEast()
}

func (r *Rocks) Load() int {
	load := 0
	for row := 0; row < len(r.slots); row += 1 {
		for col := 0; col < len(r.slots[row]); col += 1 {
			if r.slots[row][col] == 'O' {
				load += len(r.slots) - row
			}
		}
	}
	return load
}

func (r *Rocks) Clone() *Rocks {
	slots := [][]byte{}
	for _, slot := range r.slots {
		slots = append(slots, slices.Clone(slot))
	}
	return &Rocks{slots}
}

func (r *Rocks) Equal(other *Rocks) bool {
	for i := range r.slots {
		if slices.Compare(r.slots[i], other.slots[i]) != 0 {
			return false
		}
	}
	return true
}

func (r *Rocks) Dump(name string) {
	file, _ := os.Create(name)
	for _, slot := range r.slots {
		file.Write(slot)
		file.Write([]byte("\n"))
	}
}

func main() {
	file, _ := os.Open("input")
	rocks := NewRocks(file)
	rocks.TiltNorth()
	fmt.Printf("Load on north edge is %v\n", rocks.Load())
	last := []*Rocks{}
	// This is THE worst code ever LOL
	// I have a circuits final in like 13 hours so that's my excuse
	for i := 0; i < 1_000_000_000; i += 1 {
		// Tilt cycle
		rocks.TiltCycle()
		// Find how many matches exist in the previous tilt cycles
		matches := []int{}
		for i := range last {
			if last[i].Equal(rocks) {
				matches = append(matches, i)
			}
		}
		// Add ourselves to the previous cycles
		last = append(last, rocks.Clone())
		// This is very arbitrary and honestly not necessary
		if len(matches) <= 4 {
			continue
		}

		// If we have more than 4 matches calculate the spacing between them
		// Basically how long a cycle is
		delta := matches[1] - matches[0]

		// Assert that the spacing is consistent
		for i := 0; i < len(matches)-1; i += 1 {
			if matches[i+1]-matches[i] != delta {
				panic("Delta isn't consistent between matches")
			}
		}

		// Find the multiplier to get close to 1 billion - 1
		mult := (999_999_999 - matches[0]) / delta
		// If the first match plus the multiplier times the delta equals 1 billion - 1
		// we found our match (its 999 something because we're 0 indexing)
		if (matches[0] + mult*delta) == 999_999_999 {
			fmt.Printf("Load after 1 billion cycles is %v\n", rocks.Load())
			return
		}
	}
}
