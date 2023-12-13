package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day12 struct {
	day.DayInput
}

func NewDay12(inputFile string) Day12 {
	return Day12{day.DayInput(inputFile)}
}

func parseLayout(s string) []int {
	n := strings.Split(s, ",")
	result := make([]int, len(n))
	for i, j := range n {
		k, _ := strconv.Atoi(j)
		result[i] = k
	}
	return result
}

type memo struct {
	recordLen int
	cache     map[int]int
}

func newCache(recordLen, groupLen int) *memo {
	c := make(map[int]int, recordLen*groupLen)
	return &memo{recordLen, c}
}

func (m *memo) get(r, g int) (int, bool) {
	v, ok := m.cache[g*m.recordLen+r]
	return v, ok
}

func (m *memo) set(r, g, value int) {
	m.cache[g*m.recordLen+r] = value
}

func (m *memo) memoCountLayouts(record string, groups []int) int {
	if v, ok := m.get(len(record), len(groups)); ok {
		return v
	}

	v := m.countLayouts(record, groups)
	m.set(len(record), len(groups), v)
	return v
}

func hasGroupPrefix(record string, group int) bool {
	if group+1 > len(record) {
		return false
	}
	for i := 0; i < group; i++ {
		if record[i] == '.' {
			return false
		}
	}
	return record[group] != '#'
}

func (m *memo) countLayouts(record string, groups []int) int {
	if len(record) == 0 {
		if len(groups) == 0 {
			return 1
		}
		return 0
	}

	result := 0
	if record[0] != '#' {
		result += m.memoCountLayouts(record[1:], groups)
	}
	if record[0] != '.' && len(groups) > 0 && hasGroupPrefix(record, groups[0]) {
		result += m.memoCountLayouts(record[groups[0]+1:], groups[1:])
	}
	return result
}

func (d Day12) Part2() int {
	lines, _ := d.ReadLines()

	sum := 0
	for _, line := range lines {
		r1, l1, _ := strings.Cut(line, " ")
		record := strings.Join([]string{r1, r1, r1, r1, r1}, "?") + "."
		l := strings.Join([]string{l1, l1, l1, l1, l1}, ",")
		layout := parseLayout(l)

		m := newCache(len(record), len(layout))
		c := m.countLayouts(record, layout)
		sum += c
	}

	return sum
}

func (d Day12) Part1() int {
	lines, _ := d.ReadLines()

	sum := 0
	for _, line := range lines {
		record, l, _ := strings.Cut(line, " ")
		record += "."
		layout := parseLayout(l)
		m := newCache(len(record), len(layout))
		c := m.countLayouts(record, layout)
		sum += c
	}

	return sum
}

func main() {
	d := NewDay12(filepath.Join(projectpath.Root, "cmd", "day12", "input.txt"))

	day.Solve(d)
}
