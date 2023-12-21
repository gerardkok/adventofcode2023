package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay19(filepath.Join(projectpath.Root, "cmd", "day19", "example.txt"))

	want := 19114
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()
	d := NewDay19(filepath.Join(projectpath.Root, "cmd", "day19", "example.txt"))

	want := 167409079868000
	got := d.Part2()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
