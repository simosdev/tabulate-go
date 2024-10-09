package tabulate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateTableWithInstance(t *testing.T) {
	tab := New("a", "b", "c", "d")
	tab.AddRow(Row{"a": 1, "b": 2, "c": 3})
	tab.AddRow(Row{"a": 10, "b": 20, "c": 30})
	tab.AddRow(Row{"b": 2})
	tab.AddRow(Row{"d": "something"})
	tab.AddRow(Row{"d": 1.23})
	out := bytes.NewBuffer(nil)
	err := tab.Print(out)
	if err != nil {
		t.Fatalf("Print failed: %s", err)
	}

	expectedLines := []string{
		"|  a |  b |  c |         d |",
		"|----|----|----|-----------|",
		"|  1 |  2 |  3 |           |",
		"| 10 | 20 | 30 |           |",
		"|    |  2 |    |           |",
		"|    |    |    | something |",
		"|    |    |    |      1.23 |",
		"",
	}

	res := out.String()
	t.Logf("res:\n%s", res)
	lines := strings.Split(res, "\n")
	if len(lines) != len(expectedLines) {
		t.Errorf("got %d lines, want %d lines", len(lines), len(expectedLines))
		if len(lines) > len(expectedLines) {
			for i := len(expectedLines); i < len(lines); i++ {
				t.Logf("extra line [%d]: %q", i, lines[i])
			}
		}
	}
	for i, exp := range expectedLines {
		if i >= len(lines) {
			continue
		}
		if lines[i] != exp {
			t.Logf("got [%d]:  %q", i, lines[i])
			t.Logf("want [%d]: %q", i, exp)
			t.Errorf("difference on line index: %d", i)
		}
	}
}

func TestColumnMaxLengths(t *testing.T) {
	tab := New("a", "b", "c", "d")
	tab.AddRow(Row{"a": 1, "b": 200, "c": 3})
	tab.AddRow(Row{"a": 10, "b": 2, "c": 3})
	tab.AddRow(Row{"d": "something"})
	tab.AddRow(Row{"d": 1.23})
	res := tab.colMaxValueLengths()

	expected := map[string]int{"a": 2, "b": 3, "c": 1, "d": 9}
	if !cmp.Equal(expected, res) {
		t.Errorf("colMaxLengths unexpected return value")
		t.Logf("diff: %s", cmp.Diff(expected, res))
	}

}
