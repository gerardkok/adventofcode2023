package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay23(filepath.Join(projectpath.Root, "cmd", "day23", "example.txt"))

	want := 94
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay23(filepath.Join(projectpath.Root, "cmd", "day23", "example.txt"))

	want := 154
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
