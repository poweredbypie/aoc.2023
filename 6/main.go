package main

import (
	"bufio"
	"fmt"
	"math"
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

type Race struct {
	Duration int
	Record   int
}

// Part A: multiple races
func NewRaces(file *os.File) []Race {
	_, err := file.Seek(0, 0)
	PanicIf(err)

	vals := func(line string) []int {
		regex := regexp.MustCompile("[0-9]+")
		nums := []int{}
		for _, str := range regex.FindAllString(line, -1) {
			num, err := strconv.Atoi(str)
			PanicIf(err)
			nums = append(nums, num)
		}
		return nums
	}

	races := []Race{}
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "Time: ") {
			for _, num := range vals(line) {
				races = append(races, Race{Duration: num})
			}
		} else if strings.HasPrefix(line, "Distance: ") {
			for idx, num := range vals(line) {
				races[idx].Record = num
			}
		} else {
			panic("Line has unexpected content: " + line)
		}
	}

	return races
}

// Part B: one big race
func NewRace(file *os.File) Race {
	_, err := file.Seek(0, 0)
	PanicIf(err)

	race := Race{}
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		// Remove _ALL_ spaces
		line = strings.ReplaceAll(line, " ", "")
		val, err := strconv.Atoi(strings.Split(line, ":")[1])
		PanicIf(err)
		if strings.HasPrefix(line, "Time:") {
			race.Duration = val
		} else if strings.HasPrefix(line, "Distance:") {
			race.Record = val
		} else {
			panic("Line has unexpected content: " + line)
		}
	}

	return race
}

// Get the values required to reach the current record
// Small: the smaller value
// Big: the bigger value
func (r *Race) GetRecordInputs() (small float64, big float64) {
	// A race is defined by the following formula:
	// f(x) = (r.Duration - x) * x
	// To find the intercepts for f(x) == r.Record:
	// (r.Duration - x) * x == r.Record
	// Rewriting:
	// -x^2 + r.Duration * x - r.Record == 0
	// Or:
	// x^2 - r.Duration * x + r.Record == 0
	// We can plug this into the quadratic formula to solve.
	// The intercepts are:
	// (r.Duration +- sqrt((r.Duration)^2 - 4 * r.Record))) / 2
	dur := float64(r.Duration)
	rec := float64(r.Record)
	// Discriminant to reuse
	discrim := math.Sqrt((dur * dur) - 4*rec)
	small = (dur - discrim) / 2
	big = (dur + discrim) / 2
	return
}

func (r *Race) BeatCount() int {
	// A race is defined by f(x) in GetRecordInputs
	// So once we solve for the record inputs, we know the range between
	// `small` and `big` is yields all larger than the record
	small, big := r.GetRecordInputs()
	min := int(math.Ceil(small))
	max := int(math.Floor(big))
	// The range is inclusive of min, so we need to add 1 to include it.
	return (max - min + 1)
}

func PartA(races []Race) int {
	sum := 1
	for _, race := range races {
		sum *= race.BeatCount()
	}
	return sum
}

func main() {
	file, err := os.Open("input")
	PanicIf(err)
	defer file.Close()

	// Part A
	races := NewRaces(file)
	fmt.Printf("Part A sum is %v\n", PartA(races))

	// Part B
	race := NewRace(file)
	fmt.Printf("Part B single race has %v possibilities\n", race.BeatCount())
}
