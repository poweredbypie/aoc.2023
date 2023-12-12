package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type SpringInfo struct {
	groups  []int
	damaged string
}

func NewInfo(line string) SpringInfo {
	split := strings.Split(line, " ")
	damaged := split[0]
	groups := []int{}
	for _, val := range strings.Split(split[1], ",") {
		num, _ := strconv.Atoi(val)
		groups = append(groups, num)
	}
	return SpringInfo{
		groups:  groups,
		damaged: damaged,
	}
}

func Permute(str string, index int, callback func(string)) {
	if index == len(str) {
		callback(str)
		return
	}
	if str[index] == '?' {
		runes := []rune(str)
		runes[index] = '#'
		Permute(string(runes), index+1, callback)
		runes[index] = '.'
		Permute(string(runes), index+1, callback)
	} else {
		Permute(str, index+1, callback)
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

func (s *SpringInfo) GetComboCount() int {
	combos := 0
	Permute(s.damaged, 0, func(perm string) {
		if s.Validate(perm) {
			log.Printf("Validate(%v) == %v (%v) true", perm, s.damaged, s.groups)
			combos += 1
		}
	})
	return combos
}

func main() {
	file, _ := os.Open("input")
	defer file.Close()
	scan := bufio.NewScanner(file)
	sum := 0
	for scan.Scan() {
		line := scan.Text()
		info := NewInfo(line)
		count := info.GetComboCount()
		log.Printf("Count for %v: %v", line, count)
		sum += count
	}
	fmt.Printf("Sum of all combinations for all lines is %v\n", sum)
}
