package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay25b(filepath.Join(projectpath.Root, "cmd", "day25", "example.txt"))

	want := 54
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
