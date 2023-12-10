package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay06(filepath.Join(projectpath.Root, "cmd", "day06", "example.txt"))

	want := 288
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay06(filepath.Join(projectpath.Root, "cmd", "day06", "example.txt"))

	want := 71503
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
