package main

import (
	"bufio"
	"fmt"
	"os"
)

type Rocks struct {
	slots [][]byte
}

func NewRocks(file *os.File) Rocks {
	scan := bufio.NewScanner(file)
	slots := [][]byte{}
	for scan.Scan() {
		slots = append(slots, []byte(scan.Text()))
	}
	return Rocks{slots}
}

func (r *Rocks) TiltNorth() {
	tiltRock := func(row, col int) {
		for iter := row; iter >= 1; iter -= 1 {
			rock := &r.slots[iter][col]
			above := &r.slots[iter-1][col]
			if *above != '.' {
				break
			}
			// Swap
			*above, *rock = *rock, *above
		}
	}
	// Starts at 1 since row 0 is already "tilted"
	for row := 1; row < len(r.slots); row += 1 {
		for col := 0; col < len(r.slots[row]); col += 1 {
			if r.slots[row][col] == 'O' {
				tiltRock(row, col)
			}
		}
	}
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
	rocks.Dump("out")
	fmt.Printf("Load on north edge is %v\n", rocks.Load())
}
