package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExample1Part1(t *testing.T) {
	t.Parallel()
	d := NewDay20(filepath.Join(projectpath.Root, "cmd", "day20", "example1.txt"))

	want := 32000000
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExample2Part1(t *testing.T) {
	t.Parallel()
	d := NewDay20(filepath.Join(projectpath.Root, "cmd", "day20", "example2.txt"))

	want := 11687500
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
