package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func ForEachLine(file *os.File, f func(string)) {
	_, err := file.Seek(0, 0)
	PanicIf(err)

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		f(line)
	}
}

// A specific number of colored cubes.
type Draw struct {
	Color string
	Count int
}

// One "round" of each game.
type Set struct {
	Draws []Draw
}

var setRegex = regexp.MustCompile("([0-9]+) (red|green|blue)")

func NewSet(str string) Set {
	set := Set{}
	set.Draws = []Draw{}

	actions := setRegex.FindAllStringSubmatch(str, -1)
	for _, action := range actions {
		count, err := strconv.Atoi(action[1])
		PanicIf(err)
		set.Draws = append(set.Draws, Draw{
			Color: action[2],
			Count: count,
		})
	}

	return set
}

// A game; this is represented by a line in the input file
type Game struct {
	Id   int
	Sets []Set
}

var idRegex = regexp.MustCompile("Game ([0-9]+):")

func NewGame(line string) Game {
	game := Game{}
	var err error
	game.Id, err = strconv.Atoi(idRegex.FindStringSubmatch(line)[1])
	PanicIf(err)

	for _, set := range strings.Split(line, ";") {
		game.Sets = append(game.Sets, NewSet(set))
	}

	return game
}

func PartA(file *os.File) int {
	countForColor := func(draw Draw) int {
		switch draw.Color {
		case "red":
			return 12
		case "green":
			return 13
		case "blue":
			return 14
		default:
			panic("Invalid color supplied: " + draw.Color)
		}
	}

	sum := 0
	ForEachLine(file, func(line string) {
		game := NewGame(line)
		for _, set := range game.Sets {
			for _, draw := range set.Draws {
				if countForColor(draw) < draw.Count {
					return
				}
			}
		}

		sum += game.Id
	})

	return sum
}

func PartB(file *os.File) int {
	powerSum := 0
	ForEachLine(file, func(line string) {
		game := NewGame(line)
		maxRed := 0
		maxGreen := 0
		maxBlue := 0
		for _, set := range game.Sets {
			for _, draw := range set.Draws {
				var maxPtr *int
				switch draw.Color {
				case "red":
					maxPtr = &maxRed
				case "green":
					maxPtr = &maxGreen
				case "blue":
					maxPtr = &maxBlue
				}

				if draw.Count > *maxPtr {
					*maxPtr = draw.Count
				}
			}
		}

		power := maxRed * maxGreen * maxBlue
		powerSum += power
	})

	return powerSum
}

func main() {
	file, err := os.Open("input")
	PanicIf(err)
	defer file.Close()
	fmt.Printf("Sum of valid games for part A is %v\n", PartA(file))
	fmt.Printf("Powersum for part B is %v\n", PartB(file))
}
