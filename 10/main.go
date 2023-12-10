package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Coord struct {
	Row int
	Col int
}

func (c Coord) Format(f fmt.State, verb rune) {
	f.Write([]byte(fmt.Sprintf("(%v, %v)", c.Row, c.Col)))
}

func (c Coord) Equal(coord Coord) bool {
	return c.Row == coord.Row && c.Col == coord.Col
}

type Map struct {
	lines [][]byte
}

func NewMap(file *os.File) Map {
	lines := [][]byte{}
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		lines = append(lines, []byte(scan.Text()))
	}
	return Map{lines}
}

type Dir int

const (
	Up Dir = iota
	Left
	Right
	Down
)

func (d Dir) Format(f fmt.State, verb rune) {
	text := "Unknown"
	switch d {
	case Up:
		text = "Up"
	case Left:
		text = "Left"
	case Right:
		text = "Right"
	case Down:
		text = "Down"
	}
	f.Write([]byte(text))
}

func (d Dir) Flip() Dir {
	switch d {
	case Up:
		return Down
	case Left:
		return Right
	case Right:
		return Left
	case Down:
		return Up
	default:
		panic("Unexpected direction")
	}
}

func DirsFor(char byte) []Dir {
	switch char {
	case '|':
		return []Dir{Up, Down}
	case '-':
		return []Dir{Left, Right}
	case 'L':
		return []Dir{Up, Right}
	case 'J':
		return []Dir{Up, Left}
	case '7':
		return []Dir{Left, Down}
	case 'F':
		return []Dir{Right, Down}
	case '.':
		return []Dir{}
	default:
		panic("Unexpected char " + string(char))
	}
}

func (m *Map) At(coord Coord) *byte {
	return &m.lines[coord.Row][coord.Col]
}

func (m *Map) Move(coord Coord, out Dir) (Coord, Dir) {
	switch out {
	case Up:
		return Coord{coord.Row - 1, coord.Col}, out
	case Left:
		return Coord{coord.Row, coord.Col - 1}, out
	case Right:
		return Coord{coord.Row, coord.Col + 1}, out
	case Down:
		return Coord{coord.Row + 1, coord.Col}, out
	default:
		panic("Unexpected direction")
	}
}

func (m *Map) MoveDbg(coord Coord, out Dir) (Coord, Dir) {
	newVal, out := m.Move(coord, out)
	log.Printf("%v {%c} -(%v)> %v {%c}", coord, *m.At(coord), out, newVal, *m.At(newVal))
	return newVal, out
}

func (m *Map) Follow(coord Coord, in Dir) (Coord, Dir) {
	val := *m.At(coord)
	in = in.Flip()
	dirs := DirsFor(val)
	if len(dirs) != 2 {
		panic("Too many or too little directions")
	}
	if dirs[0] == in {
		return m.MoveDbg(coord, dirs[1])
	} else {
		return m.MoveDbg(coord, dirs[0])
	}
}

func (m *Map) Loop(start Coord, dir Dir, dists map[Coord]int) {
	for curr, dir, dist := start, dir.Flip(), 1; ; dist += 1 {
		curr, dir = m.Follow(curr, dir)
		// This needs to happen after the first check
		if curr.Equal(start) {
			log.Printf("Stopped at distance %v", dist)
			break
		}
		lastDist, ok := dists[curr]
		if ok {
			log.Printf("Existing distance is %v, new is %v", lastDist, dist)
		}
		if lastDist > dist || !ok {
			dists[curr] = dist
		}
		log.Printf("Distance %v for coord %v", dist, curr)
	}
}

func (m *Map) ConnectTrack(start Coord) byte {
	all := []Dir{Up, Left, Right, Down}
	dirs := []Dir{}
	for _, dir := range all {
		around, _ := m.Move(start, dir)
		val := *m.At(around)
		dirsAround := DirsFor(val)
		flipped := dir.Flip()
		if dirsAround[0] == flipped || dirsAround[1] == flipped {
			dirs = append(dirs, dir)
		}
	}
	if len(dirs) != 2 {
		panic("Too many or too little directions for ConnectTrack")
	}
	slices.Sort(dirs)
	var char byte
	switch {
	case dirs[0] == Up && dirs[1] == Left:
		char = 'J'
	case dirs[0] == Up && dirs[1] == Right:
		char = 'L'
	case dirs[0] == Up && dirs[1] == Down:
		char = '|'
	case dirs[0] == Left && dirs[1] == Right:
		char = '-'
	case dirs[0] == Left && dirs[1] == Down:
		char = '7'
	case dirs[0] == Right && dirs[1] == Down:
		char = 'F'
	default:
		panic("Couldn't match dirs")
	}
	*m.At(start) = char
	return char
}

func (m *Map) Start() Coord {
	coord := Coord{0, 0}
	for coord.Row = 0; coord.Row < len(m.lines); coord.Row += 1 {
		for coord.Col = 0; coord.Col < len(m.lines[coord.Row]); coord.Col += 1 {
			if *m.At(coord) == 'S' {
				return coord
			}
		}
	}
	panic("Couldn't find start of loop")
}

func main() {
	file, err := os.Open("input")
	PanicIf(err)
	m := NewMap(file)
	start := m.Start()
	// Connect start with actual nodes it connects to
	startChar := m.ConnectTrack(start)
	dirs := DirsFor(startChar)
	if len(dirs) != 2 {
		panic("Too many / too little directions")
	}
	dists := make(map[Coord]int)
	m.Loop(start, dirs[0], dists)
	m.Loop(start, dirs[1], dists)
	max := 0
	for _, val := range dists {
		if val > max {
			max = val
		}
	}
	fmt.Printf("Max value for distance is %v\n", max)
}
