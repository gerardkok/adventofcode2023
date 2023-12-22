package main

import (
	"path/filepath"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day20 struct {
	day.DayInput
}

func NewDay20(inputFile string) Day20 {
	return Day20{day.DayInput(inputFile)}
}

type module struct {
	mtype                 moduleType
	on                    bool
	inputs                map[string]bool
	sources, destinations []string
}

type moduleType int

const (
	dummy moduleType = iota
	flipflop
	conjunction
	broadcaster
)

type pulse struct {
	source, destination string
	high                bool
}

type machine map[string]module

func parseSource(src string) (string, moduleType) {
	switch src[0] {
	case '%':
		return src[1:], flipflop
	case '&':
		return src[1:], conjunction
	default:
		return src, broadcaster
	}
}

func parseLines(lines []string) machine {
	result := make(machine, len(lines))
	for _, line := range lines {
		src, dest, _ := strings.Cut(line, " -> ")
		dests := strings.Split(dest, ", ")
		name, mtype := parseSource(src)
		result[name] = module{
			mtype:        mtype,
			on:           false,
			inputs:       make(map[string]bool),
			sources:      make([]string, 0),
			destinations: dests,
		}
	}

	//connect
	for k, v := range result {
		for _, dest := range v.destinations {
			s := result[dest]
			s.sources = append(s.sources, k)
			result[dest] = s
		}
	}

	return result
}

func (m moduleType) String() string {
	switch m {
	case flipflop:
		return "flip-flop"
	case conjunction:
		return "conjunction"
	case broadcaster:
		return "broadcaster"
	default:
		return "dummy"
	}
}

func (p pulse) String() string {
	if p.high {
		return "-high->"
	}
	return "-low->"
}

func (m module) allHigh() bool {
	for _, src := range m.sources {
		if !m.inputs[src] {
			return false
		}
	}
	return true
}

func (m *module) processFlipflop(p pulse) []pulse {
	result := make([]pulse, 0)
	if p.high {
		return result
	}

	m.on = !m.on
	for _, d := range m.destinations {
		send := pulse{p.destination, d, m.on}
		result = append(result, send)
	}

	return result
}

func (m *module) processConjunction(p pulse) []pulse {
	result := make([]pulse, 0)
	m.inputs[p.source] = p.high
	allHigh := m.allHigh()
	for _, d := range m.destinations {
		send := pulse{p.destination, d, !allHigh}
		result = append(result, send)
	}
	return result
}

func (m *module) processBroadcaster(p pulse) []pulse {
	result := make([]pulse, 0)
	for _, d := range m.destinations {
		send := pulse{p.destination, d, p.high}
		result = append(result, send)
	}
	return result
}

func (m *module) processPulse(p pulse) []pulse {
	switch m.mtype {
	case flipflop:
		return m.processFlipflop(p)
	case conjunction:
		return m.processConjunction(p)
	case broadcaster:
		return m.processBroadcaster(p)
	default:
		return []pulse{}
	}
}

func (m machine) pushButton() []pulse {
	done := make([]pulse, 0)
	todo := make([]pulse, 0)
	todo = append(todo, pulse{"button", "broadcaster", false})

	for len(todo) > 0 {
		pulse := todo[0]
		todo = todo[1:]

		mod := m[pulse.destination]
		pulses := mod.processPulse(pulse)
		todo = append(todo, pulses...)
		m[pulse.destination] = mod
		done = append(done, pulse)
	}

	return done
}

func countHigh(pulses []pulse) int {
	result := 0
	for _, p := range pulses {
		if p.high {
			result++
		}
	}
	return result
}

func (d Day20) Part1() int {
	lines, _ := d.ReadLines()
	machine := parseLines(lines)

	nLow, nHigh := 0, 0

	for i := 0; i < 1000; i++ {
		pulses := machine.pushButton()
		h := countHigh(pulses)
		nLow, nHigh = nLow+len(pulses)-h, nHigh+h
	}

	return nLow * nHigh
}

func (m machine) findRxSources() map[string]int {
	for _, v := range m {
		for _, d := range v.destinations {
			if d == "rx" {
				result := make(map[string]int)
				for _, s := range v.sources {
					result[s] = 0
				}
				return result
			}
		}
	}
	return nil
}

func lowPulseToRx(sources map[string]int) int {
	result := 1
	for _, v := range sources {
		result *= v
	}
	return result
}

func (d Day20) Part2() int {
	lines, _ := d.ReadLines()
	machine := parseLines(lines)
	rxSources := machine.findRxSources()

	i := 0
	for lowPulseToRx(rxSources) == 0 {
		i++
		pulses := machine.pushButton()
		for _, p := range pulses {
			if !p.high {
				continue
			}
			if _, ok := rxSources[p.source]; !ok {
				continue
			}
			rxSources[p.source] = i
		}
	}

	return lowPulseToRx(rxSources)
}

func main() {
	d := NewDay20(filepath.Join(projectpath.Root, "cmd", "day20", "input.txt"))

	day.Solve(d)
}
