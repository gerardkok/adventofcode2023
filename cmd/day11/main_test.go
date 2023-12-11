package main

import (
	"path/filepath"
	"testing"

	"adventofcode23/internal/projectpath"
)

func TestExamplePart1(t *testing.T) {
	t.Parallel()
	d := NewDay11(filepath.Join(projectpath.Root, "cmd", "day11", "example.txt"), 2, 0)

	want := 374
	got := d.Part1()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestExamplePart2(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expansion int
		want      int
	}{
		{
			expansion: 10,
			want:      1030,
		},
		{
			expansion: 100,
			want:      8410,
		},
	}

	for i, tc := range testCases {
		d := NewDay11(filepath.Join(projectpath.Root, "cmd", "day11", "example.txt"), 0, tc.expansion)
		got := d.Part2()
		if tc.want != got {
			t.Errorf("test %d: want %d, got %d", i, tc.want, got)
		}
	}
}
