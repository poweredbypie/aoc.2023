package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
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

func DirsFor(char byte) (Dir, Dir) {
	switch char {
	case '|':
		return Up, Down
	case '-':
		return Left, Right
	case 'L':
		return Up, Right
	case 'J':
		return Up, Left
	case '7':
		return Left, Down
	case 'F':
		return Right, Down
	default:
		panic("Unexpected char " + string(char))
	}
}

type Map struct {
	lines  [][]byte
	rows   int
	cols   int
	border map[Coord]bool
}

func MapFromFile(file *os.File) Map {
	lines := [][]byte{}
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		lines = append(lines, []byte(scan.Text()))
	}
	return NewMap(lines)
}

func NewMap(bytes [][]byte) Map {
	return Map{
		lines:  bytes,
		rows:   len(bytes),
		cols:   len(bytes[0]),
		border: make(map[Coord]bool),
	}
}

func (m *Map) At(coord Coord) *byte {
	return &m.lines[coord.Row][coord.Col]
}

// Stuff for part A (follow the track around)
func (m *Map) Move(coord Coord, out Dir) (Coord, Dir) {
	// Assumption: loop does not lead into border
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

func (m *Map) MoveDebug(coord Coord, out Dir) (Coord, Dir) {
	newVal, out := m.Move(coord, out)
	log.Printf("%v {%c} -(%v)> %v {%c}", coord, *m.At(coord), out, newVal, *m.At(newVal))
	return newVal, out
}

func (m *Map) Follow(coord Coord, from Dir) (Coord, Dir) {
	val := *m.At(coord)
	from = from.Flip()
	in, out := DirsFor(val)
	if in == from {
		return m.MoveDebug(coord, out)
	} else {
		return m.MoveDebug(coord, in)
	}
}

func (m *Map) Loop(start Coord, dir Dir) int {
	for curr, dir, dist := start, dir.Flip(), 1; ; dist += 1 {
		m.border[curr] = true
		curr, dir = m.Follow(curr, dir)
		// This needs to happen after the first check
		if curr.Equal(start) {
			log.Printf("Stopped at distance %v", dist)
			return dist / 2
		}
	}
}

func (m *Map) ConnectTrack(start Coord) byte {
	all := []Dir{Up, Left, Right, Down}
	dirs := []Dir{}
	for _, dir := range all {
		around, _ := m.Move(start, dir)
		val := *m.At(around)
		in, out := DirsFor(val)
		flip := dir.Flip()
		if in == flip || out == flip {
			dirs = append(dirs, dir)
		}
	}
	if len(dirs) != 2 {
		panic("Too many or too little directions for ConnectTrack")
	}
	// "Sort" so the order is consistent
	slices.Sort(dirs)
	first, second := dirs[0], dirs[1]
	var char byte
	switch {
	case first == Up && second == Left:
		char = 'J'
	case first == Up && second == Right:
		char = 'L'
	case first == Up && second == Down:
		char = '|'
	case first == Left && second == Right:
		char = '-'
	case first == Left && second == Down:
		char = '7'
	case first == Right && second == Down:
		char = 'F'
	default:
		panic("Couldn't match dirs")
	}
	*m.At(start) = char
	return char
}

func (m *Map) Start() Coord {
	c := Coord{}
	for c.Row = 0; c.Row < m.rows; c.Row += 1 {
		for c.Col = 0; c.Col < m.cols; c.Col += 1 {
			if *m.At(c) == 'S' {
				return c
			}
		}
	}
	panic("Couldn't find start of loop")
}

// Stuff for part B (fill the outside and find leftover inside
func (m *Map) Dump(filename string) {
	file, err := os.Create(filename)
	PanicIf(err)
	for _, line := range m.lines {
		file.Write(line)
		file.Write([]byte("\n"))
	}
}

const UpDownLarge = `
.|.
.|.
.|.
`
const LeftRightLarge = `
...
---
...
`
const UpRightLarge = `
.|.
.L-
...
`
const UpLeftLarge = `
.|.
-J.
...
`
const LeftDownLarge = `
...
-7.
.|.
`
const RightDownLarge = `
...
.F-
.|.
`

func GetLarge(char byte) [][]byte {
	str := ""
	switch char {
	case '|':
		str = UpDownLarge
	case '-':
		str = LeftRightLarge
	case 'L':
		str = UpRightLarge
	case 'J':
		str = UpLeftLarge
	case '7':
		str = LeftDownLarge
	case 'F':
		str = RightDownLarge
	default:
		panic("Unexpected char " + string(char))
	}
	bytes := [][]byte{}
	for _, str := range strings.Split(str, "\n") {
		if str != "" {
			bytes = append(bytes, []byte(str))
		}
	}
	log.Printf("Got %v byte arrays in %v", len(bytes), bytes)
	return bytes
}

func (m *Map) WriteAt(coord Coord, bytes [][]byte) {
	for row := 0; row < len(bytes); row += 1 {
		for col := 0; col < len(bytes[row]); col += 1 {
			curr := Coord{row + coord.Row, col + coord.Col}
			*m.At(curr) = bytes[row][col]
		}
	}
}

// Expand each square into 9 squares
// This allows a fill operation to pass through gaps that wouldn't be possible otherwise
func (m *Map) MakeLarge() Map {
	lines := [][]byte{}
	for row := 0; row < m.rows*3; row += 1 {
		line := make([]byte, m.cols*3)
		for col := range line {
			line[col] = '.'
		}
		lines = append(lines, line)
	}
	large := NewMap(lines)
	for coord := range m.border {
		log.Printf("Processing coord %v", coord)
		bytes := GetLarge(*m.At(coord))
		// Get new coord by "large"
		coord.Row *= 3
		coord.Col *= 3
		large.WriteAt(coord, bytes)
	}
	return large
}

// Convert a 9 square big map back into a small map
func (m *Map) MakeSmall() Map {
	lines := [][]byte{}
	for row := 0; row < m.rows/3; row += 1 {
		line := make([]byte, m.cols/3)
		for col := range line {
			line[col] = '.'
		}
		lines = append(lines, line)
	}
	small := NewMap(lines)
	c := Coord{}
	for c.Row = 0; c.Row < m.rows; c.Row += 3 {
		for c.Col = 0; c.Col < m.cols; c.Col += 3 {
			middle := c
			middle.Row += 1
			middle.Col += 1
			smallC := c
			smallC.Row /= 3
			smallC.Col /= 3
			// Pick the square at the center of each 3x3
			*small.At(smallC) = *m.At(middle)
		}
	}
	return small
}

func (m *Map) FindAround(coord Coord, char byte) bool {
	// Only checking as such:
	// .*.
	// *S*
	// .*.
	// where * is checked.
	// The small board needs to check diagonals, but the big board does not because it's expanded.
	// Above
	if coord.Row > 0 {
		curr := coord
		curr.Row -= 1
		if *m.At(curr) == char {
			return true
		}
	}
	// Left
	if coord.Col > 0 {
		curr := coord
		curr.Col -= 1
		if *m.At(curr) == char {
			return true
		}
	}
	// Right
	if coord.Col < m.cols-1 {
		curr := coord
		curr.Col += 1
		if *m.At(curr) == char {
			return true
		}
	}
	// Below
	if coord.Row < m.rows-1 {
		curr := coord
		curr.Row += 1
		if *m.At(curr) == char {
			return true
		}
	}
	return false
}

func (m *Map) FillOutside() {
	isBorder := func(c byte) bool {
		return c == '|' || c == '-' || c == 'L' || c == 'J' || c == '7' || c == 'F'
	}
	loop := func(c Coord, nextRow func(int) int, nextCol func(int) int) int {
		sum := 0
		// Assumption: start coord is outside
		*m.At(c) = 'O'
		rowStart, colStart := c.Row, c.Col
		for c.Row = rowStart; c.Row >= 0 && c.Row < m.rows; c.Row = nextRow(c.Row) {
			for c.Col = colStart; c.Col >= 0 && c.Col < m.cols; c.Col = nextCol(c.Col) {
				if *m.At(c) != 'O' && !isBorder(*m.At(c)) && m.FindAround(c, 'O') {
					*m.At(c) = 'O'
					sum += 1
				}
			}
		}
		return sum
	}
	inc := func(val int) int { return val + 1 }
	dec := func(val int) int { return val - 1 }
	// We have to do this 4 separate times for each direction, iterating until we get no new fills
	// Otherwise we can't fill all cavities
	for {
		sum := 0
		// First: L -> R, U -> D
		sum += loop(Coord{0, 0}, inc, inc)
		// Second: R -> L, U -> D
		sum += loop(Coord{0, m.cols - 1}, inc, dec)
		// Third: L -> R, D -> U
		sum += loop(Coord{m.rows - 1, 0}, dec, inc)
		// Fourth: R -> L, D -> U
		sum += loop(Coord{m.rows - 1, m.cols - 1}, dec, dec)
		if sum == 0 {
			break
		}
	}
}

func (m *Map) Count(char byte) int {
	c := Coord{}
	count := 0
	for c.Row = 0; c.Row < m.rows; c.Row += 1 {
		for c.Col = 0; c.Col < m.cols; c.Col += 1 {
			if *m.At(c) == char {
				count += 1
			}
		}
	}
	return count
}

func main() {
	file, err := os.Open("input")
	defer file.Close()
	PanicIf(err)
	m := MapFromFile(file)
	start := m.Start()
	// Connect start with actual nodes it connects to
	startChar := m.ConnectTrack(start)
	dir, _ := DirsFor(startChar)
	fmt.Printf("Max value for distance is %v\n", m.Loop(start, dir))
	large := m.MakeLarge()
	large.FillOutside()
	// For visualization
	large.Dump("big")
	small := large.MakeSmall()
	small.Dump("small")
	fmt.Printf("Number of unfilled values is %v\n", small.Count('.'))
}
