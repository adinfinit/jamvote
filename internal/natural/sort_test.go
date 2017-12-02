package natural

import (
	"reflect"
	"sort"
	"testing"
)

func TestLess(t *testing.T) {
	var items = []string{
		"A200",
		"A40",
		"A3",
		"A5",
		"A",
		"Alpha 100",
		"Alpha 10",
	}

	sort.Slice(items, func(i, k int) bool {
		return Less(items[i], items[k])
	})

	var expected = []string{
		"A",
		"A3",
		"A5",
		"A40",
		"A200",
		"Alpha 10",
		"Alpha 100",
	}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Got:\n%v\n%v", items, expected)
	}
}
