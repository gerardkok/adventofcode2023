package main

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type handType int

type Day07 struct {
	day.DayInput
}

func NewDay07(inputFile string) Day07 {
	return Day07{day.DayInput(inputFile)}
}

const (
	highCard handType = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind

	faceOrder1 = "23456789TJQKA"
	faceOrder2 = "J23456789TQKA"
)

var (
	highCardPattern     = []int{1, 1, 1, 1, 1}
	onePairPattern      = []int{1, 1, 1, 2}
	twoPairPattern      = []int{1, 2, 2}
	threeOfAKindPattern = []int{1, 1, 3}
	fullHousePattern    = []int{2, 3}
	fourOfAKindPattern  = []int{1, 4}
	fiveOfAKindPattern  = []int{5}
)

type handBid struct {
	hand     string
	bid      int
	handType handType
}

func equalPatterns(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func pattern1(hand string) []int {
	cardValues := make(map[rune]int)
	for _, c := range hand {
		cardValues[c]++
	}
	pattern := make([]int, 0)
	for _, v := range cardValues {
		pattern = append(pattern, v)
	}
	sort.Ints(pattern)
	return pattern
}

func pattern2(hand string) []int {
	cardValues := make(map[rune]int)
	for _, c := range hand {
		cardValues[c]++
	}
	pattern := make([]int, 0)
	for k, v := range cardValues {
		if k != 'J' {
			pattern = append(pattern, v)
		}
	}
	sort.Ints(pattern)
	if len(pattern) == 0 {
		// JJJJJ
		return fiveOfAKindPattern
	}
	pattern[len(pattern)-1] += cardValues['J']
	return pattern
}

func getHandType(hand string, patternFn func(string) []int) handType {
	pattern := patternFn(hand)
	switch {
	case equalPatterns(pattern, highCardPattern):
		return highCard
	case equalPatterns(pattern, onePairPattern):
		return onePair
	case equalPatterns(pattern, twoPairPattern):
		return twoPair
	case equalPatterns(pattern, threeOfAKindPattern):
		return threeOfAKind
	case equalPatterns(pattern, fullHousePattern):
		return fullHouse
	case equalPatterns(pattern, fourOfAKindPattern):
		return fourOfAKind
	case equalPatterns(pattern, fiveOfAKindPattern):
		return fiveOfAKind
	}
	return highCard
}

func newHandBid(hand string, bid int, patternFn func(string) []int) handBid {
	handType := getHandType(hand, patternFn)
	return handBid{hand, bid, handType}
}

func less(a, b handBid, faceOrder string) bool {
	if a.handType < b.handType {
		return true
	}
	if a.handType > b.handType {
		return false
	}

	for k := 0; k < 5; k++ {
		faceValueI := strings.IndexByte(faceOrder, a.hand[k])
		faceValueJ := strings.IndexByte(faceOrder, b.hand[k])
		if faceValueI < faceValueJ {
			return true
		}
		if faceValueI > faceValueJ {
			return false
		}
	}
	return false
}

func winnings(handBids []handBid) int {
	result := 0
	for i, h := range handBids {
		w := (i + 1) * h.bid
		result += w
	}
	return result
}

func (d Day07) Part1() int {
	input, _ := d.ReadLines()
	handBids := make([]handBid, len(input))
	for i, line := range input {
		hand, b, _ := strings.Cut(line, " ")
		bid, _ := strconv.Atoi(b)
		handBids[i] = newHandBid(hand, bid, pattern1)
	}

	sort.Slice(handBids, func(i, j int) bool {
		return less(handBids[i], handBids[j], faceOrder1)
	})

	return winnings(handBids)
}

func (d Day07) Part2() int {
	input, _ := d.ReadLines()
	handBids := make([]handBid, len(input))
	for i, line := range input {
		hand, b, _ := strings.Cut(line, " ")
		bid, _ := strconv.Atoi(b)
		handBids[i] = newHandBid(hand, bid, pattern2)
	}

	sort.Slice(handBids, func(i, j int) bool {
		return less(handBids[i], handBids[j], faceOrder2)
	})

	return winnings(handBids)
}

func main() {
	d := NewDay07(filepath.Join(projectpath.Root, "cmd", "day07", "input.txt"))

	day.Solve(d)
}
