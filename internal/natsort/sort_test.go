// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package natsort

import (
	"fmt"
	"slices"
	"testing"
)

var tests = []struct {
	a, b string
	want int
}{
	{"", "", 0},
	{"a", "", 1},
	{"abc", "abc", 0},
	{"ab", "abc", -1},
	{"x", "ab", 1},
	{"x", "a", 1},
	{"b", "x", -1},

	{"a", "aa", -1},
	{"a", "a1", -1},
	{"ab", "a1", 1},
	{"a", "0", 1},
	{"A", "0", 1},
	{"file1.txt", "file2.txt", -1},
	{"file10.txt", "file2.txt", 1},
	{"file1000.txt", "file2.txt", 1},
	{"file0001.txt", "file2.txt", -1},
	{"file00a.txt", "file000a.txt", -1},
	{"a10", "a010", -1},
	{"a1b2", "a01b3", -1},
	{"file_1.txt", "file_10.txt", -1},
	{"file1.txt", "file1.txt", 0},
	{"fileA.txt", "fileB.txt", -1},
	{"file1A.txt", "file1B.txt", -1},
	{"Uint8x16", "Uint32x8", -1},
	{"Uint32x16", "Uint32x8", 1},
	{"Uint10000000000000000000000", "Uint20000000000000000000000", -1},
	{"Uint10000000000000000000000abc", "Uint10000000000000000000000abc", 0},
	{"a1a1a1a1a1a1a1a1a1a1a1", "a1a1a1a1a1a1a1a1a1a1a1", 0},
	{"a1a1a1a1a1a1a1a1a1a1a10", "a1a1a1a1a1a1a1a1a1a1a1", 1},
}

func TestCompare(t *testing.T) {
	for _, tt := range tests {
		if got := Compare(tt.a, tt.b); got != tt.want {
			t.Errorf("Compare(%q, %q) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
		if got := Compare(tt.b, tt.a); got != -tt.want {
			t.Errorf("Compare(%q, %q) = %d; want %d", tt.b, tt.a, got, -tt.want)
		}
	}
}

func TestSliceSort(t *testing.T) {
	types := []string{"Uint32x16", "Uint16x32", "Unit8x64", "Uint64x8"}
	want := []string{"Unit8x64", "Uint16x32", "Uint32x16", "Uint64x8"}
	slices.SortFunc(types, Compare)
	if slices.Equal(types, want) {
		t.Errorf("types = %v; want %v", types, want)
	}
}

func BenchmarkCompare(b *testing.B) {
	for i, test := range tests {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				Compare(test.a, test.b)
			}
		})
	}
}

func FuzzTransitivity(f *testing.F) {
	f.Add("", "", "")
	f.Fuzz(func(t *testing.T, a string, b string, c string) {
		ab := Compare(a, b)
		bc := Compare(b, c)
		ca := Compare(c, a)

		// when the total is 3 or -3, it means that there is a cycle in comparison.
		if tot := ab + bc + ca; tot == 3 || tot == -3 {
			t.Errorf("Compare(%q, %q) = %d", a, b, ab)
			t.Errorf("Compare(%q, %q) = %d", b, c, bc)
			t.Errorf("Compare(%q, %q) = %d", c, a, ca)
		}
	})
}
