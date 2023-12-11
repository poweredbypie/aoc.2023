package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Schematic struct {
	Lines      []string
	Rows, Cols int
}

func Digit(char byte) bool {
	return unicode.IsDigit(rune(char))
}

type Coord struct {
	Row, Col int
}

func (c Coord) Format(f fmt.State, _ rune) {
	f.Write([]byte(fmt.Sprintf("(%v, %v)", c.Row, c.Col)))
}

func (s *Schematic) At(coord Coord) byte {
	return s.Lines[coord.Row][coord.Col]
}

type SchematicIter struct {
	schem *Schematic
	curr  Coord
}

func (s *Schematic) NewIterator() SchematicIter {
	return SchematicIter{
		// Initially use  -1 for column since you need to call Next() first
		s, Coord{0, -1},
	}
}

func (s *SchematicIter) Next() bool {
	s.curr.Col += 1
	if s.curr.Col >= s.schem.Cols {
		s.curr.Row += 1
		s.curr.Col = 0
		if s.curr.Row >= s.schem.Rows {
			return false
		}
	}

	return true
}

func (s *SchematicIter) Coord() Coord {
	return s.curr
}

func (s *Schematic) CheckAround(coord Coord, compare func(char byte, coord Coord) bool) bool {
	compareWrap := func(rowVal, colVal int) bool {
		coord := Coord{rowVal, colVal}
		return compare(s.At(coord), coord)
	}
	check := func(rowVal int) bool {
		if rowVal >= 0 && rowVal < s.Rows {
			// Check specified row
			if coord.Col > 0 {
				// Check specified row, to the left
				if compareWrap(rowVal, coord.Col-1) {
					return true
				}
			}
			// Check specified row, same coord.Column
			if compareWrap(rowVal, coord.Col) {
				return true
			}
			if coord.Col < s.Cols-1 {
				// Check specified row, to the right
				if compareWrap(rowVal, coord.Col+1) {
					return true
				}
			}
		}
		return false
	}
	return check(coord.Row-1) || check(coord.Row) || check(coord.Row+1)
}

// For part A
func (s *Schematic) FindValidNumbers() []int {
	keep := false
	numStr := ""
	nums := []int{}

	isSymbol := func(char byte, _ Coord) bool {
		if Digit(char) {
			return false
		} else {
			return char != '.'
		}
	}

	iter := s.NewIterator()
	for iter.Next() {
		char := s.At(iter.Coord())

		if Digit(char) {
			// Still in a number (or starting a new one)
			numStr += string(char)
			if s.CheckAround(iter.Coord(), isSymbol) {
				keep = true
			}
		} else {
			// Number finished, check states
			if keep {
				num, err := strconv.Atoi(numStr)
				PanicIf(err)
				nums = append(nums, num)
			}
			// Reset states for next number
			keep = false
			numStr = ""
		}
	}

	return nums
}

// Find all stars (for part B)
func (s *Schematic) FindStars() []Coord {
	coords := []Coord{}

	iter := s.NewIterator()
	for iter.Next() {
		if s.At(iter.Coord()) == '*' {
			coords = append(coords, iter.Coord())
		}
	}

	return coords
}

// Get the leftmost digit for this number
func (s *Schematic) ParentDigit(coord Coord) Coord {
	iter := coord
	for ; iter.Col >= 0; iter.Col -= 1 {
		if !Digit(s.At(iter)) {
			iter.Col += 1
			return iter
		}
	}
	// Cap column at 0. Can't be -1
	iter.Col = 0
	return iter
}

func (s *Schematic) ParentDigitToNum(parent Coord) int {
	findRight := func() Coord {
		iter := parent
		for ; iter.Col < s.Cols; iter.Col += 1 {
			if !Digit(s.At(iter)) {
				iter.Col -= 1
				return iter
			}
		}
		// Cap column at s.Cols - 1
		iter.Col -= 1
		return iter
	}
	right := findRight()
	slice := s.Lines[parent.Row][parent.Col : right.Col+1]
	num, err := strconv.Atoi(slice)
	PanicIf(err)
	return num
}

func NewSchematic(file *os.File) Schematic {
	_, err := file.Seek(0, 0)
	PanicIf(err)
	scan := bufio.NewScanner(file)
	schem := Schematic{}
	for scan.Scan() {
		schem.Lines = append(schem.Lines, scan.Text())
	}

	schem.Rows = len(schem.Lines)
	schem.Cols = len(schem.Lines[0])

	for idx, line := range schem.Lines {
		if len(line) != schem.Cols {
			panic(fmt.Sprintf("Line %v has differing length from first line (%v vs. %v)", idx, schem.Cols, len(line)))
		}
	}

	return schem
}

func PartA(schem Schematic) int {
	sum := 0
	for _, num := range schem.FindValidNumbers() {
		sum += num
	}
	return sum
}

func PartB(schem Schematic) uint64 {
	sum := uint64(0)
	for _, star := range schem.FindStars() {
		parentMap := make(map[Coord]bool)
		schem.CheckAround(star, func(char byte, coord Coord) bool {
			if Digit(char) {
				// Digits can belong to the same number, so use this to determine it
				parentMap[schem.ParentDigit(coord)] = true
			}
			return false
		})
		if len(parentMap) != 2 {
			continue
		}
		// Convert the map to a slice
		parents := []Coord{}
		for parent := range parentMap {
			parents = append(parents, parent)
		}
		ratio := schem.ParentDigitToNum(parents[0]) * schem.ParentDigitToNum(parents[1])
		sum += uint64(ratio)
	}

	return sum
}

func main() {
	file, err := os.Open("input")
	PanicIf(err)
	schem := NewSchematic(file)
	file.Close()
	fmt.Printf("Part A sum is %v\n", PartA(schem))
	fmt.Printf("Part B sum is %v\n", PartB(schem))
}
