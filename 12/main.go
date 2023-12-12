package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type CacheKey struct {
	bufLen int
	group  int
}

type SpringInfo struct {
	groups  []int
	damaged string
	cache   map[CacheKey]int
}

func NewInfo(line string) *SpringInfo {
	split := strings.Split(line, " ")
	// Beginning and ending dots are basically useless to us
	damaged := strings.Trim(split[0], ".")
	groups := []int{}
	for _, val := range strings.Split(split[1], ",") {
		num, _ := strconv.Atoi(val)
		groups = append(groups, num)
	}
	return &SpringInfo{
		groups:  groups,
		damaged: damaged,
		cache:   make(map[CacheKey]int),
	}
}

func NewUnfoldedInfo(line string) *SpringInfo {
	split := strings.Split(line, " ")
	// Disgusting
	join5 := func(str string, sep string) string {
		arr := []string{str, str, str, str, str}
		return strings.Join(arr, sep)
	}
	damaged := join5(split[0], "?")
	groups := []int{}
	// EW
	for _, val := range strings.Split(join5(split[1], ","), ",") {
		num, _ := strconv.Atoi(val)
		groups = append(groups, num)
	}
	return &SpringInfo{
		groups:  groups,
		damaged: damaged,
		cache:   make(map[CacheKey]int),
	}
}

func (s *SpringInfo) CountRecur(buf []rune, grpIdx int) int {
	// Algorithm:
	// Find the first "slot" that can fit the current group
	// Place the group in the slot
	// Move on to the next group by calling ourselves again with the next group and a buf slice
	// Find the next slot in the group and try again until there are no more slots

	// Base case: we filled all groups
	if grpIdx >= len(s.groups) {
		// If we can find a '#' between here and the end, we failed
		// Otherwise, we succeeded! Return 1 to add to the count
		if slices.Contains(buf, '#') {
			return 0
		} else {
			return 1
		}
	}

	// Base case: trying to fill a group with an empty buffer
	// We can't place it anywhere so just return with 0 count
	if len(buf) == 0 {
		return 0
	}

	key := CacheKey{len(buf), grpIdx}

	// Base case: if the value is cached return that instead
	if val, ok := s.cache[key]; ok {
		return val
	}

	// Recursive case: find the next slot to put the next group in before continuing
	groupLen := s.groups[grpIdx]
	// Helper functions
	// Find length of slot
	findLen := func(idx int) int {
		iter := idx
		for iter = idx; iter < len(buf) && buf[iter] != '.'; iter += 1 {
		}
		return iter - idx
	}
	// Fill buffer with group
	fillSlot := func(idx int) {
		for iter := idx; iter < idx+groupLen; iter += 1 {
			buf[iter] = '#'
		}
	}
	// Skip separators (non-broken springs)
	skipDots := func(idx int) int {
		iter := idx
		for iter = idx; iter < len(buf) && buf[iter] == '.'; iter += 1 {
		}
		return iter
	}

	slotLen := 0
	// Total combinations found
	sum := 0
	// Loop until we find a slot that can fit the group, or until we reach the end
	for idx := skipDots(0); idx < len(buf); idx = skipDots(idx + slotLen) {
		slotLen = findLen(idx)
		if slotLen < groupLen {
			continue
		}
		// We found a slot with enough space, now loop through all combinations until we can't anymore
		for off := 0; off <= slotLen-groupLen; off += 1 {
			// Calculate the offset for the next group
			end := idx + off + groupLen
			// If the separator cannot be replaced with a '.', or a '#' is before the offset,
			// we can't use this.
			// Only check the separator if we're not at the end of the buffer
			if (end < len(buf) && buf[end] == '#') || slices.Contains(buf[:idx+off], '#') {
				continue
			}

			// To store old values
			old := make([]rune, groupLen)

			// Get the existing values in the slot to replace later
			copy(old, buf[idx+off:])
			fillSlot(idx + off)
			// This just makes it so the next group doesn't need to check for a '.' at the beginning
			nextOff := min(end+1, len(buf))
			sum += s.CountRecur(buf[nextOff:], grpIdx+1)
			// Paste back in
			copy(buf[idx+off:], old)
		}
	}
	// Cache the value we calculated
	s.cache[key] = sum
	return sum
}

func (s *SpringInfo) GetComboCount() int {
	return s.CountRecur([]rune(s.damaged), 0)
}

func main() {
	file, _ := os.Open("input")
	defer file.Close()
	scan := bufio.NewScanner(file)
	sum := 0
	unfoldSum := 0
	for scan.Scan() {
		line := scan.Text()
		sum += NewInfo(line).GetComboCount()
		unfoldSum += NewUnfoldedInfo(line).GetComboCount()
	}
	fmt.Printf("Sum of all combinations for all lines is %v\n", sum)
	fmt.Printf("Sum of all combinations for all unfolded lines is %v\n", unfoldSum)
}
