package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	strsToNums := func(strs []string) []int {
		nums := []int{}
		for _, str := range strs {
			if str == "" {
				continue
			}
			num, err := strconv.Atoi(str)
			PanicIf(err)
			nums = append(nums, num)
		}
		return nums
	}
	file, err := os.Open("input")
	PanicIf(err)
	scan := bufio.NewScanner(file)

	part1Sum := 0
	cardIdx := 1
	copies := [300]int{}
	for scan.Scan() {
		nums := strings.Split(scan.Text(), ":")[1]
		split := strings.Split(nums, "|")
		given := strsToNums(strings.Split(split[0], " "))
		mine := strsToNums(strings.Split(split[1], " "))

		part1Score := 0
		part2Score := 0

		for _, num := range given {
			if slices.Contains(mine, num) {
				part2Score += 1
				if part1Score == 0 {
					part1Score = 1
				} else {
					part1Score *= 2
				}
			}
		}

		for idx := 0; idx < part2Score; idx += 1 {
			copies[idx+1+cardIdx] += 1 + copies[cardIdx]
		}

		part1Sum += part1Score
		cardIdx += 1
	}

	part2Sum := 0
	for idx := 1; idx < cardIdx; idx += 1 {
		part2Sum += 1 + copies[idx]
	}

	fmt.Printf("Part 1 sum is %v\n", part1Sum)
	fmt.Printf("Part 2 sum is %v\n", part2Sum)
}
