package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

var debug = log.New(os.Stderr, "DEBUG: ", 0)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type Remap struct {
	SrcStart int
	DstStart int
	Length   int
}

func (r *Remap) GetMapped(src int) (int, bool) {
	if src >= r.SrcStart && src < r.SrcStart+r.Length {
		return r.DstStart + (src - r.SrcStart), true
	}
	return 0, false
}

type Map struct {
	From string
	To   string
	// This is sorted by SrcStart
	Remaps []*Remap
}

func (m *Map) GetMapping(src int) int {
	for _, remap := range m.Remaps {
		mapping, mapped := remap.GetMapped(src)
		if mapped {
			return mapping
		}
	}
	return src
}

type SeedRange struct {
	Start  int
	Length int
}

// Since part B has huge ranges, we can't have all seeds pre-allocated in an array
type Seeds struct {
	nexted   bool
	currIdx  int
	currSeed int
	Ranges   []SeedRange
}

func NewSeeds() Seeds {
	return Seeds{
		nexted:  false,
		currIdx: 0,
	}
}

func getSeedList(line string) []string {
	seedList := []string{}

	if strings.HasPrefix(line, "seeds: ") {
		seedList = strings.Split(
			strings.Trim(
				strings.Split(line, ":")[1],
				" ",
			), " ",
		)
	}

	return seedList
}

func NewList(line string) Seeds {
	seeds := NewSeeds()
	list := getSeedList(line)
	for _, seed := range list {
		id, err := strconv.Atoi(seed)
		PanicIf(err)
		seeds.Ranges = append(seeds.Ranges, SeedRange{
			Start:  id,
			Length: 1,
		})
		debug.Printf("Adding seed ID %v", id)
	}

	return seeds
}

func NewRanges(line string) Seeds {
	seeds := NewSeeds()
	list := getSeedList(line)
	for idx := 0; idx < len(list); idx += 2 {
		// Start and length pairs
		start, err := strconv.Atoi(list[idx])
		PanicIf(err)
		length, err := strconv.Atoi(list[idx+1])
		PanicIf(err)
		seeds.Ranges = append(seeds.Ranges, SeedRange{start, length})
		debug.Printf("Adding seed range %v to %v", start, start+length-1)
	}

	return seeds
}

func (s *Seeds) Count() int {
	count := 0
	for _, currRange := range s.Ranges {
		count += currRange.Length
	}
	return count
}

func (s *Seeds) Next() bool {
	currRange := s.Ranges[s.currIdx]
	if !s.nexted {
		s.nexted = true
		s.currSeed = currRange.Start
		return true
	}

	s.currSeed += 1
	if s.currSeed >= currRange.Start+currRange.Length {
		s.currIdx += 1
		if s.currIdx >= len(s.Ranges) {
			return false
		}
		currRange = s.Ranges[s.currIdx]
		debug.Printf("Starting new range (%v to %v)", currRange.Start, currRange.Start+currRange.Length-1)
		s.currSeed = currRange.Start
	}

	return true
}

func (s *Seeds) Seed() int {
	if !s.nexted {
		panic("Didn't call s.Next() first")
	}
	return s.currSeed
}

func GetRemaps(scan *bufio.Scanner) []*Remap {
	remaps := []*Remap{}

	for scan.Scan() {
		line := scan.Text()
		if line == "" {
			break
		}

		mappings := strings.Split(line, " ")
		dest, err := strconv.Atoi(mappings[0])
		PanicIf(err)
		src, err := strconv.Atoi(mappings[1])
		PanicIf(err)
		length, err := strconv.Atoi(mappings[2])
		PanicIf(err)

		debug.Printf("Found remap from %v to %v (length %v)", src, dest, length)

		remaps = append(remaps, &Remap{
			SrcStart: src,
			DstStart: dest,
			Length:   length,
		})
	}

	slices.SortFunc(remaps, func(one, two *Remap) int {
		if one.SrcStart < two.SrcStart {
			return -1
		} else if one.SrcStart > two.SrcStart {
			return 1
		} else {
			return 0
		}
	})

	return remaps
}

type Maps map[string]*Map

func GetMaps(scan *bufio.Scanner) Maps {
	maps := make(Maps)

	for scan.Scan() {
		line := scan.Text()
		if strings.Contains(line, "map:") {
			mapName := strings.Split(strings.TrimSuffix(line, " map:"), "-")
			from := mapName[0]
			to := mapName[2]
			debug.Printf("Found map from %v to %v", from, to)
			maps[from] = &Map{
				From:   from,
				To:     to,
				Remaps: GetRemaps(scan),
			}
		}
	}

	return maps
}

func (m *Maps) Follow(start string, src int) int {
	next := start
	val := src
	for mapping := (*m)[next]; mapping != nil; mapping = (*m)[next] {
		next = mapping.To
		newVal := mapping.GetMapping(val)
		val = newVal
	}
	// This is too noisy for part B
	// debug.Printf("Mapping %v (%v) to %v (%v)", src, start, val, next)

	return val
}

func (m *Maps) GetMin(seeds *Seeds) int {
	min := int(^uint(0) >> 1)
	for seeds.Next() {
		val := m.Follow("seed", seeds.Seed())
		if val < min {
			min = val
		}
	}

	return min
}

func main() {
	file, err := os.Open("input")
	PanicIf(err)
	defer file.Close()
	scan := bufio.NewScanner(file)

	// Read the first line for both types of seed listings
	scan.Scan()
	firstLine := scan.Text()

	partA := NewList(firstLine)
	// This takes like 10 or 15 minutes to run LOL don't do this
	// At least space complexity is super flat
	partB := NewRanges(firstLine)

	maps := GetMaps(scan)
	fmt.Printf("Part B has %v seeds\n", partB.Count())
	fmt.Printf("Part A: minimum mapped value is %v\n", maps.GetMin(&partA))
	fmt.Printf("Part B: minimum mapped value is %v\n", maps.GetMin(&partB))
}
