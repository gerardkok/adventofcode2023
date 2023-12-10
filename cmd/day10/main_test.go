package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExample1Part2(t *testing.T) {
	t.Parallel()
	p := NewPuzzle10(filepath.Join(projectpath.Root, "cmd", "day10", "example1-part2.txt"))

	want := 4
	got := p.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part2(t *testing.T) {
	t.Parallel()
	p := NewPuzzle10(filepath.Join(projectpath.Root, "cmd", "day10", "example2-part2.txt"))

	want := 8
	got := p.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample3Part2(t *testing.T) {
	t.Parallel()
	p := NewPuzzle10(filepath.Join(projectpath.Root, "cmd", "day10", "example3-part2.txt"))

	want := 10
	got := p.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
