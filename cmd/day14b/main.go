package main

import (
	"os"
	"strings"

	"github.com/cespare/xxhash/v2"

	"adventofcode23/internal/day"
)

type Day14b struct {
	day.DayInput
}

type platform struct {
	nRows, nColumns int
	spots           []byte
}

func NewDay14b(inputFile string) Day14b {
	return Day14b{day.DayInput(inputFile)}
}

func (p *platform) tiltNorth() {
	north := make([]int, p.nColumns)

	for r := 0; r < p.nRows; r++ {
		for c := 0; c < p.nColumns; c++ {
			i := r*p.nColumns + c
			switch p.spots[i] {
			case 'O':
				j := north[c]*p.nColumns + c
				p.spots[i], p.spots[j] = p.spots[j], p.spots[i]
				north[c]++
			case '#':
				north[c] = r + 1
			}
		}
	}
}

func (p *platform) tiltWest() {
	west := make([]int, p.nRows)

	for c := 0; c < p.nColumns; c++ {
		for r := 0; r < p.nRows; r++ {
			i := r*p.nColumns + c
			switch p.spots[i] {
			case 'O':
				j := r*p.nColumns + west[r]
				p.spots[i], p.spots[j] = p.spots[j], p.spots[i]
				west[r]++
			case '#':
				west[r] = c + 1
			}
		}
	}
}

func (p *platform) tiltSouth() {
	south := make([]int, p.nColumns)
	for i := range south {
		south[i] = p.nRows - 1
	}

	for r := p.nRows - 1; r >= 0; r-- {
		for c := 0; c < p.nColumns; c++ {
			i := r*p.nColumns + c
			switch p.spots[i] {
			case 'O':
				j := south[c]*p.nColumns + c
				p.spots[i], p.spots[j] = p.spots[j], p.spots[i]
				south[c]--
			case '#':
				south[c] = r - 1
			}
		}
	}
}

func (p *platform) tiltEast() {
	east := make([]int, p.nRows)
	for i := range east {
		east[i] = p.nColumns - 1
	}

	for c := p.nColumns - 1; c >= 0; c-- {
		for r := 0; r < p.nRows; r++ {
			i := r*p.nColumns + c
			switch p.spots[i] {
			case 'O':
				j := r*p.nColumns + east[r]
				p.spots[i], p.spots[j] = p.spots[j], p.spots[i]
				east[r]--
			case '#':
				east[r] = c - 1
			}
		}
	}
}

func (d Day14b) Part1() int {
	lines, _ := d.ReadLines()

	p := makePlatform(lines)
	p.tiltNorth()

	return p.load()
}

func (p platform) load() int {
	result := 0
	for r := 0; r < p.nRows; r++ {
		l := p.nRows - r
		for c := r * p.nColumns; c < (r+1)*p.nColumns; c++ {
			if p.spots[c] == 'O' {
				result += l
			}
		}
	}
	return result
}

func (p *platform) cycle() {
	p.tiltNorth()
	p.tiltWest()
	p.tiltSouth()
	p.tiltEast()
}

func (p *platform) detectLoop() (int, []int) {
	loads := make([]int, 0)
	seen := make(map[uint64]int)
	for {
		xxh := xxhash.Sum64(p.spots)
		if index, ok := seen[xxh]; ok {
			return index, loads[index:]
		}
		seen[xxh] = len(loads)
		loads = append(loads, p.load())
		p.cycle()
	}
}

func makePlatform(lines []string) platform {
	nRows := len(lines)
	nColumns := len(lines[0])
	spots := strings.Join(lines, "")
	return platform{nRows, nColumns, []byte(spots)}
}

func (d Day14b) Part2() int {
	lines, _ := d.ReadLines()

	p := makePlatform(lines)

	s, e := p.detectLoop()

	last := (1000000000 - s) % len(e)

	return e[last]
}

func main() {
	d := NewDay14b(os.Args[1])

	day.Solve(d)
}
