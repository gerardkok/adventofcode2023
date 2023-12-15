package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day15 struct {
	day.DayInput
}

type lens struct {
	label    string
	focalLen int
	loc      int
}

type box struct {
	lenses map[string]*lens
	order  []*lens
}

func NewDay15(inputFile string) Day15 {
	return Day15{day.DayInput(inputFile)}
}

func newBox() box {
	lenses := make(map[string]*lens)
	order := make([]*lens, 0)
	return box{lenses, order}
}

func (b *box) add(label string, focalLen int) {
	if v, ok := b.lenses[label]; ok {
		// lens already present
		v.focalLen = focalLen
		return
	}
	loc := len(b.order)
	l := lens{label, focalLen, loc}
	b.lenses[label] = &l
	b.order = append(b.order, &l)
}

func (b *box) delete(label string) {
	if _, ok := b.lenses[label]; !ok {
		// lens already not present
		return
	}
	loc := b.lenses[label].loc
	b.order = append(b.order[:loc], b.order[loc+1:]...)
	for i := loc; i < len(b.order); i++ {
		b.order[i].loc--
	}
	delete(b.lenses, label)
}

func (b box) power() int {
	sum := 0
	for i := range b.order {
		l := b.order[i]
		sum += (i + 1) * l.focalLen
	}
	return sum
}

func HASH(s string) int {
	result := 0
	for _, c := range s {
		result += int(c)
		result *= 17
		result %= 256
	}
	return result
}

func (d Day15) Part1() int {
	lines, _ := d.ReadLines()

	steps := strings.Split(lines[0], ",")

	sum := 0
	for _, step := range steps {
		sum += HASH(step)
	}
	return sum
}

func perform(step string, boxes []box) {
	if label, f, ok := strings.Cut(step, "="); ok {
		boxID := HASH(label)
		focalLen, _ := strconv.Atoi(f)
		boxes[boxID].add(label, focalLen)
	} else {
		label := strings.TrimSuffix(step, "-")
		boxID := HASH(label)
		boxes[boxID].delete(label)
	}
}

func (d Day15) Part2() int {
	lines, _ := d.ReadLines()

	boxes := make([]box, 256)
	for i := range boxes {
		boxes[i] = newBox()
	}

	steps := strings.Split(lines[0], ",")
	for _, step := range steps {
		perform(step, boxes)
	}

	sum := 0
	for i, b := range boxes {
		sum += (i + 1) * b.power()
	}
	return sum
}

func main() {
	d := NewDay15(filepath.Join(projectpath.Root, "cmd", "day15", "input.txt"))

	day.Solve(d)
}
