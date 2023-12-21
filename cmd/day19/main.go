package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"
)

type Day19 struct {
	day.DayInput
}

func NewDay19(inputFile string) Day19 {
	return Day19{day.DayInput(inputFile)}
}

const (
	accepted  = "A"
	rejected  = "R"
	minRating = 1
	maxRating = 4000
)

type part map[byte]int

type interval struct {
	min, max int // min inclusive, max exclusive
}

type intervalMap map[byte]interval

type rule struct {
	category         byte
	target           string
	accepts, rejects interval
}

type workflows map[string][]rule

func parseInput(lines []string) (workflows, []part) {
	result := make([][]string, 2)
	result = append(result, make([]string, 0))

	i := 0
	for _, line := range lines {
		if len(line) == 0 {
			i++
			result = append(result, make([]string, 0))
			continue
		}

		result[i] = append(result[i], line)
	}

	w := parseWorkflows(result[0])
	r := parseParts(result[1])

	return w, r
}

func parseCondition(condition, target string) rule {
	category := condition[0]
	threshold, _ := strconv.Atoi(condition[2:])
	accepts := interval{threshold + 1, maxRating}
	rejects := interval{minRating, threshold}
	if condition[1] == byte('<') {
		accepts = interval{minRating, threshold - 1}
		rejects = interval{threshold, maxRating}
	}

	return rule{
		category: category,
		target:   target,
		accepts:  accepts,
		rejects:  rejects,
	}
}

func parseRule(r string) rule {
	condition, target, ok := strings.Cut(r, ":")
	if !ok {
		// just a target
		return rule{
			target:  r,
			accepts: interval{0, 0}, // accept everything
		}
	}
	return parseCondition(condition, target)
}

func parseRules(rules string) []rule {
	s := strings.Split(rules, ",")
	result := make([]rule, len(s))

	for i, r := range s {
		result[i] = parseRule(r)
	}
	return result
}

func parseWorkflows(w []string) workflows {
	result := make(workflows, len(w))
	for _, workflow := range w {
		name, r, _ := strings.Cut(workflow[:len(workflow)-1], "{") // cut off }
		result[name] = parseRules(r)
	}
	return result
}

func parsePart(rating string) part {
	result := make(part, 4)
	split := strings.Split(rating[1:len(rating)-1], ",") // cut off { and }

	for _, r := range split {
		category := r[0]
		value, _ := strconv.Atoi(r[2:])
		result[category] = value
	}
	return result
}

func parseParts(ratings []string) []part {
	result := make([]part, len(ratings))
	for i, rating := range ratings {
		result[i] = parsePart(rating)
	}
	return result
}

func (r rule) applies(v int) bool {
	return v >= r.accepts.min && v <= r.accepts.max
}

func (w workflows) accept(p part, workflow string) bool {
	if workflow == accepted || workflow == rejected {
		return workflow == accepted
	}

	rules := w[workflow]
	for _, r := range rules {
		if r.applies(p[r.category]) {
			return w.accept(p, r.target)
		}
	}
	return false
}

func (p part) sum() int {
	result := 0
	for _, v := range p {
		result += v
	}
	return result
}

func (i interval) length() int {
	if i.min > i.max {
		return 0
	}
	return i.max - i.min + 1
}

func (i intervalMap) countSolutions() int {
	result := 1
	for _, j := range i {
		result *= j.length()
	}
	return result
}

func intersect(a, b interval) interval {
	return interval{max(a.min, b.min), min(a.max, b.max)}
}

func (i intervalMap) clone() intervalMap {
	result := make(intervalMap, 4)
	for k, v := range i {
		result[k] = v
	}
	return result
}

func (w workflows) countSolutions(intervals intervalMap, workflow string) int {
	switch workflow {
	case accepted:
		return intervals.countSolutions()
	case rejected:
		return 0
	}

	count := 0

	rules := w[workflow]
	for _, r := range rules {
		cat := intervals[r.category]
		newIntervals := intervals.clone()
		newIntervals[r.category] = intersect(cat, r.accepts)
		intervals[r.category] = intersect(cat, r.rejects)
		count += w.countSolutions(newIntervals, r.target)
	}

	return count
}

func (d Day19) Part1() int {
	lines, _ := d.ReadLines()
	workflows, ratings := parseInput(lines)

	sum := 0

	for _, rating := range ratings {
		if workflows.accept(rating, "in") {
			sum += rating.sum()
		}
	}

	return sum
}

func (d Day19) Part2() int {
	lines, _ := d.ReadLines()
	workflows, _ := parseInput(lines)

	intervals := intervalMap{
		'x': {minRating, maxRating},
		'm': {minRating, maxRating},
		'a': {minRating, maxRating},
		's': {minRating, maxRating},
	}

	return workflows.countSolutions(intervals, "in")
}

func main() {
	d := NewDay19(filepath.Join(projectpath.Root, "cmd", "day19", "input.txt"))

	day.Solve(d)
}
