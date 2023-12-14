package main

import (
	"path/filepath"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day14 struct {
	day.DayInput
}

type platform [][]byte

func NewDay14(inputFile string) Day14 {
	return Day14{day.DayInput(inputFile)}
}

func (p *platform) tilt() {
	north := make([]int, len((*p)[0]))

	for i, r := range *p {
		for j := range r {
			switch (*p)[i][j] {
			case 'O':
				(*p)[i][j], (*p)[north[j]][j] = (*p)[north[j]][j], (*p)[i][j]
				north[j]++
			case '#':
				north[j] = i + 1
			}
		}
	}
}

func (p *platform) rotate() {
	// reverse rows
	for i, j := 0, len(*p)-1; i < j; i, j = i+1, j-1 {
		(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
	}

	// transpose
	for i := 0; i < len(*p); i++ {
		for j := 0; j < i; j++ {
			(*p)[i][j], (*p)[j][i] = (*p)[j][i], (*p)[i][j]
		}
	}
}

func (d Day14) Part1() int {
	lines, _ := d.ReadLines()

	p := makePlatform(lines)
	p.tilt()

	return p.load()
}

func (p platform) load() int {
	result := 0
	for i, row := range p {
		for _, c := range row {
			if c == 'O' {
				result += len(p) - i
			}
		}
	}
	return result
}

func (p *platform) cycle() {
	for i := 0; i < 4; i++ {
		p.tilt()
		p.rotate()
	}
}

func equal(p, q platform) bool {
	for i, row := range p {
		for j, c := range row {
			if q[i][j] != c {
				return false
			}
		}
	}
	return true
}

func (p platform) in(l []platform) int {
	for i, q := range l {
		if equal(p, q) {
			return i
		}
	}
	return -1
}

func (p platform) copy() platform {
	result := make(platform, len(p))
	for i := range p {
		result[i] = make([]byte, len(p[i]))
		copy(result[i], p[i])
	}
	return result
}

func (p *platform) findLoop() (int, []platform) {
	seen := make([]platform, 0)
	for {
		q := p.copy()
		seen = append(seen, q)
		p.cycle()
		i := p.in(seen)
		if i > -1 {
			return i, seen[i:]
		}
	}
}

func makePlatform(lines []string) platform {
	p := make(platform, len(lines))
	for i, line := range lines {
		p[i] = []byte(line)
	}
	return p
}

func (d Day14) Part2() int {
	lines, _ := d.ReadLines()

	p := makePlatform(lines)

	s, e := p.findLoop()

	last := (1000000000 - s) % len(e)

	return e[last].load()
}

func main() {
	d := NewDay14(filepath.Join(projectpath.Root, "cmd", "day14", "input.txt"))

	day.Solve(d)
}
