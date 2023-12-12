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

type CacheKey struct {
	bufLen int
	group  int
}

type SpringInfo struct {
	groups  []int
	damaged string
	combos  int
	cache   map[CacheKey]int
}

func NewInfo(line string) SpringInfo {
	split := strings.Split(line, " ")
	// Beginning and ending dots are basically useless to us
	damaged := strings.Trim(split[0], ".")
	groups := []int{}
	for _, val := range strings.Split(split[1], ",") {
		num, _ := strconv.Atoi(val)
		groups = append(groups, num)
	}
	return SpringInfo{
		groups:  groups,
		damaged: damaged,
		combos:  0,
		cache:   make(map[CacheKey]int),
	}
}

func NewUnfoldedInfo(line string) SpringInfo {
	split := strings.Split(line, " ")
	// Disgusting
	join5 := func(str string, sep string) string {
		arr := []string{str, str, str, str, str}
		return strings.Join(arr, sep)
	}
	damaged := join5(split[0], "?")
	groups := []int{}
	for _, val := range strings.Split(join5(split[1], ","), ",") {
		num, _ := strconv.Atoi(val)
		groups = append(groups, num)
	}
	return SpringInfo{
		groups:  groups,
		damaged: damaged,
		combos:  0,
		cache:   make(map[CacheKey]int),
	}
}

func (s *SpringInfo) Validate(str string) bool {
	index := 0
	for _, str := range strings.Split(str, ".") {
		// Remove all empty strings
		if str == "" {
			continue
		}
		// Too many groups in the string
		if index >= len(s.groups) {
			return false
		}
		// Group i doesn't match the length
		if len(str) != s.groups[index] {
			return false
		}
		index += 1
	}
	return index == len(s.groups)
}

func (s *SpringInfo) CountRecur(buf []rune, grpIdx int, whole []rune) int {
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
			// log.Printf("Failed with %v", string(whole))
			return 0
		} else {
			// log.Printf("Succeeded with %v (buf '%v')", string(whole), string(buf))
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
	// To store old values when we replace
	old := make([]rune, groupLen)
	// Total combinations found
	sum := 0
	// Find length of slot
	findLen := func(idx int) int {
		iter := idx
		for buf[iter] == '#' || buf[iter] == '?' {
			iter += 1
			if iter >= len(buf) {
				// Reached the end so stop increasing
				break
			}
		}
		return iter - idx
	}
	// Fill buffer with group
	fillSlot := func(idx int) {
		for iter := idx; iter < idx+groupLen; iter += 1 {
			buf[iter] = '#'
		}
	}

	bufIdx := 0
	slotLen := 0

	// Try filling the slot for every slot that we can find that can house this group
	slotLoop := func() {
		for off := 0; off <= slotLen-groupLen; off += 1 {
			// Calculate the offset for the next group
			end := bufIdx + off + groupLen
			// This just makes it so the next group doesn't need to check for a '.' at the beginning
			nextOff := min(end+1, len(buf))
			// If the separator cannot be replaced with a '.', we can't use this
			// Only check this if we're not at the end of the buffer
			if end < len(buf) && buf[end] == '#' {
				// log.Printf("Buf '%v' at %v is #, skipping", string(buf), nextOff-1)
				// log.Printf("For group %v of '%v'", grpIdx, string(whole))
				continue
			}
			// Check that our group is precisely groupLen long
			if slices.Contains(buf[:bufIdx+off], '#') {
				// log.Printf("Group %v is overcommitting, skipping", grpIdx)
				// log.Printf("For %v", string(whole))
				continue
			}
			// Get the existing values in the slot to replace later
			copy(old, buf[bufIdx+off:])
			fillSlot(bufIdx + off)
			sum += s.CountRecur(buf[nextOff:], grpIdx+1, whole)
			// Paste back in
			copy(buf[bufIdx+off:], old)
		}
	}

	// Loop until we find a slot that can fit the group, or until we reach the end
	for {
		// Skip leading dots, we can't place a group there
		for buf[bufIdx] == '.' {
			bufIdx += 1
			// We can't place a group anymore, we've reached the end of the buffer
			if bufIdx >= len(buf) {
				s.cache[key] = sum
				return sum
			}
		}
		slotLen = findLen(bufIdx)
		if slotLen >= groupLen {
			// We found a slot with enough space, now loop through all combinations until we can't anymore
			slotLoop()
		}
		// Skip the slot if we can't fit the group in there
		bufIdx += slotLen
		if bufIdx >= len(buf) {
			// Again, if we can't place a group anymore, we've reached the end of the buffer
			s.cache[key] = sum
			return sum
		}
	}
}

func (s *SpringInfo) GetComboCount() int {
	s.combos = 0
	// log.Printf("Calling count recur")
	runes := []rune(s.damaged)
	s.combos = s.CountRecur(runes, 0, runes)
	return s.combos
}

func main() {
	file, _ := os.Open("input")
	defer file.Close()
	scan := bufio.NewScanner(file)
	sumA := 0
	sumB := 0
	for scan.Scan() {
		line := scan.Text()
		infoA := NewInfo(line)
		infoB := NewUnfoldedInfo(line)
		log.Printf("Running for line %v", line)
		sumA += infoA.GetComboCount()
		sumB += infoB.GetComboCount()
	}
	fmt.Printf("Sum of all combinations for all lines is %v\n", sumA)
	fmt.Printf("Sum of all combinations for all unfolded lines is %v\n", sumB)
}
