package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"adventofcode23/internal/day"
	"adventofcode23/internal/projectpath"

	"golang.org/x/exp/maps"
)

type Day23 struct {
	day.DayInput
}

type area struct {
	neighbours map[tile][]tile
	start, end tile
}

type tile struct {
	row, column int
}

type move struct {
	dRow, dColumn int
}

type edge struct {
	to       tile
	distance int
}

type graph struct {
	edges      map[tile][]edge
	start, end tile
}

var slopes = map[byte]move{
	'^': {-1, 0},
	'>': {0, 1},
	'v': {1, 0},
	'<': {0, -1},
}

func NewDay23(inputFile string) Day23 {
	return Day23{day.DayInput(inputFile)}
}

func parseTiles(input []string) [][]byte {
	tiles := make([][]byte, len(input)+2)
	tiles[0] = []byte(strings.Repeat("#", len(input[0])+2))
	tiles[len(input)+1] = tiles[0]
	for i, line := range input {
		tiles[i+1] = []byte("#" + line + "#")
	}
	return tiles
}

func makeArea(tiles [][]byte, moves func(byte) []move) area {
	neighbours := make(map[tile][]tile)
	var start, end tile

	for r, row := range tiles {
		for c, ch := range row {
			t := tile{r, c}
			neighbours[t] = make([]tile, 0)
			for _, m := range moves(ch) {
				if tiles[r+m.dRow][c+m.dColumn] == '#' {
					continue
				}
				neighbours[t] = append(neighbours[t], tile{r + m.dRow, c + m.dColumn})
			}
			if ch == '.' {
				switch r {
				case 1:
					start = t
				case len(tiles) - 2:
					end = t
				}
			}
		}
	}
	return area{neighbours, start, end}
}

func (a area) intersections() []tile {
	result := []tile{a.start, a.end}
	for k, v := range a.neighbours {
		if len(v) > 2 {
			result = append(result, k)
		}
	}
	return result
}

func (a area) neighbourIntersections(e edge, intersections []tile, exclude int, seen map[tile]struct{}) []edge {
	for i := range intersections {
		if i == exclude {
			continue
		}
		if intersections[i] == e.to {
			return []edge{e}
		}
	}

	seen[e.to] = struct{}{}
	result := make([]edge, 0)
	for _, n := range a.neighbours[e.to] {
		if _, ok := seen[n]; ok {
			continue
		}
		edges := a.neighbourIntersections(edge{n, e.distance + 1}, intersections, exclude, seen)
		result = append(result, edges...)
	}
	return result
}

func (a area) makeGraph() graph {
	intersections := a.intersections()
	edges := make(map[tile][]edge, len(intersections))

	for i, intersection := range intersections {
		seen := make(map[tile]struct{})
		edges[intersection] = a.neighbourIntersections(edge{intersection, 0}, intersections, i, seen)
	}
	return graph{edges, a.start, a.end}
}

func (g graph) bfs(e edge, end tile, seen map[tile]struct{}) []int {
	if e.to == end {
		return []int{e.distance}
	}

	result := make([]int, 0)
	for _, n := range g.edges[e.to] {
		if _, ok := seen[n.to]; ok {
			continue
		}

		seen[e.to] = struct{}{}
		distances := g.bfs(edge{n.to, n.distance + e.distance}, end, seen)
		delete(seen, e.to)
		result = append(result, distances...)
	}

	return result
}

func printEdgeSlice(v []edge) string {
	n := make([]string, len(v))
	for i, w := range v {
		n[i] = fmt.Sprintf("%s: %d", w.to, w.distance)
	}
	return strings.Join(n, ", ")
}

func (g graph) String() string {
	b := make([]string, len(g.edges))
	j := 0
	for k, v := range g.edges {
		b[j] = fmt.Sprintf("%s: [%s]", k, printEdgeSlice(v))
		j++
	}
	edges := strings.Join(b, "\n")

	return fmt.Sprintf("edges:\n%s\nstart: %s, end: %s\n", edges, g.start, g.end)
}

func (g graph) maxDistance() int {
	seen := make(map[tile]struct{})

	distances := g.bfs(edge{g.start, 0}, g.end, seen)

	result := 0
	for _, d := range distances {
		result = max(d, result)
	}
	return result
}

func (t tile) String() string {
	return fmt.Sprintf("(%d,%d)", t.row, t.column)
}

func printTileSlice(s []tile) string {
	n := make([]string, len(s))
	for i, w := range s {
		n[i] = w.String()
	}
	return strings.Join(n, ",")
}

func (a area) String() string {
	b := make([]string, len(a.neighbours))
	j := 0
	for k, v := range a.neighbours {
		b[j] = fmt.Sprintf("%s: [%s]", k, printTileSlice(v))
		j++
	}
	adjacencies := strings.Join(b, "\n")

	return fmt.Sprintf("adjancencies:\n%s\nstart: %s, end: %s\n", adjacencies, a.start, a.end)
}

func (d Day23) Part1() int {
	lines, _ := d.ReadLines()
	tiles := parseTiles(lines)
	area := makeArea(tiles, func(ch byte) []move {
		switch ch {
		case '^', '>', 'v', '<':
			return []move{slopes[ch]}
		case '.':
			return maps.Values(slopes)
		default:
			return nil
		}
	})

	graph := area.makeGraph()

	return graph.maxDistance()
}

func (d Day23) Part2() int {
	lines, _ := d.ReadLines()
	tiles := parseTiles(lines)
	area := makeArea(tiles, func(ch byte) []move {
		switch ch {
		case '.', '^', '>', 'v', '<':
			return maps.Values(slopes)
		default:
			return nil
		}
	})

	graph := area.makeGraph()

	return graph.maxDistance()
}

func main() {
	d := NewDay23(filepath.Join(projectpath.Root, "cmd", "day23", "input.txt"))

	day.Solve(d)
}
