package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func partA(file *os.File) (int, error) {
	regex, err := regexp.Compile("[1-9]")
	if err != nil {
		return -1, errors.New("Couldn't compile regex: " + err.Error())
	}
	if _, err := file.Seek(0, 0); err != nil {
		return -1, errors.New("Couldn't seek to beginning of file: " + err.Error())
	}

	sum := 0

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		found := regex.FindAllString(line, -1)
		if found == nil {
			return -1, errors.New("Couldn't find any numbers in the line!")
		}
		str := found[0] + found[len(found)-1]
		num, err := strconv.Atoi(str)
		if err != nil {
			return -1, errors.New("Couldn't convert string to number: " + err.Error())
		}
		sum = sum + num
	}
	return sum, nil
}

func partB(file *os.File) (int, error) {
	regex, err := regexp.Compile("[1-9]|one|two|three|four|five|six|seven|eight|nine")
	if err != nil {
		return -1, errors.New("Couldn't compile regex: " + err.Error())
	}

	if _, err := file.Seek(0, 0); err != nil {
		return -1, errors.New("Couldn't seek to beginning of file: " + err.Error())
	}

	parse := func(str string) int {
		if num, err := strconv.Atoi(str); err == nil {
			return num
		}
		switch str {
		case "one":
			return 1
		case "two":
			return 2
		case "three":
			return 3
		case "four":
			return 4
		case "five":
			return 5
		case "six":
			return 6
		case "seven":
			return 7
		case "eight":
			return 8
		case "nine":
			return 9
		default:
			panic("Passed an invalid string to parse() func: " + str)
		}
	}

	sum := 0
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		found := regex.FindAllString(line, -1)
		if found == nil {
			return -1, errors.New("Couldn't find any numbers in line!")
		}
		num := parse(found[0])*10 + parse(found[len(found)-1])
		sum = sum + num
	}

	return sum, nil
}

func main() {
	file, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if sum, err := partA(file); err != nil {
		panic(err)
	} else {
		fmt.Printf("The sum for part A is %v\n", sum)
	}
	// Go doesn't support overlapping regexes so this is bad (see the Python impl)
	// Ex: oneight doesn't emit '18'
	if sum, err := partB(file); err != nil {
		panic(err)
	} else {
		fmt.Printf("The sum for part B is %v\n", sum)
	}
}
