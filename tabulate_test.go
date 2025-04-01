package tabulate

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTabulate(t *testing.T) {
	res := Tabulate(
		nil,
		Row{"a": 1, "b": 2, "c": 3},
		Row{"a": 10, "b": 20, "c": 30},
		Row{"d": 1.23},
	)

	expectedLines := []string{
		"|  a |  b |  c |    d |",
		"|----|----|----|------|",
		"|  1 |  2 |  3 |      |",
		"| 10 | 20 | 30 |      |",
		"|    |    |    | 1.23 |",
		"",
	}

	assertRows(t, expectedLines, res)
}

func TestCreateTableWithInstance(t *testing.T) {
	tab := New([]string{"a", "bbb", "c", "d"})
	tab.AddRow(Row{"a": 1, "bbb": 2, "c": 3})
	tab.AddRow(Row{"a": 10, "bbb": 20, "c": 30})
	tab.AddRow(Row{"bbb": 2})
	tab.AddRow(Row{"a": []int{1, 2, 3}})
	tab.AddRow(Row{"a": map[string]int{"h": 1, "i": 2}})
	tab.AddRow(Row{"d": "something"})
	tab.AddRow(Row{"d": 1.23})
	tab.AddRow(Row{"d": func() string { return "cb-val" }})
	out := bytes.NewBuffer(nil)
	err := tab.Print(out)
	if err != nil {
		t.Fatalf("Print failed: %s", err)
	}

	expectedLines := []string{
		"|            a | bbb |  c |         d |",
		"|--------------|-----|----|-----------|",
		"|            1 |   2 |  3 |           |",
		"|           10 |  20 | 30 |           |",
		"|              |   2 |    |           |",
		"|      [1 2 3] |     |    |           |",
		"| map[h:1 i:2] |     |    |           |",
		"|              |     |    | something |",
		"|              |     |    |      1.23 |",
		"|              |     |    |    cb-val |",
		"",
	}

	res := out.String()
	assertRows(t, expectedLines, res)
}

func TestAddRowWithAutoAddColumns(t *testing.T) {
	tab := New(nil)
	tab.AddRow(Row{"a": 1})
	tab.AddRow(Row{"a": 10, "c": 30, "b": 20})
	tab.AddRow(Row{"b": 2})

	out := bytes.NewBuffer(nil)
	err := tab.Print(out)
	if err != nil {
		t.Fatalf("Print failed: %s", err)
	}

	expectedLines := []string{
		"|  a |  b |  c |",
		"|----|----|----|",
		"|  1 |    |    |",
		"| 10 | 20 | 30 |",
		"|    |  2 |    |",
		"",
	}

	res := out.String()
	assertRows(t, expectedLines, res)
}

func TestAddRowWithExtraColumns(t *testing.T) {
	tab := New([]string{"a", "b", "c"})
	row1 := Row{"a": 1}
	row1["d"] = 3
	tab.AddRow(row1)
	tab.AddRow(Row{"a": 10, "c": 30, "b": 20})
	row3 := Row{"b": 2}
	row3["d"] = 15
	tab.AddRow(row3)

	out := bytes.NewBuffer(nil)
	err := tab.Print(out)
	if err != nil {
		t.Fatalf("Print failed: %s", err)
	}

	expectedLines := []string{
		"|  a |  b |  c |  d |",
		"|----|----|----|----|",
		"|  1 |    |    |  3 |",
		"| 10 | 20 | 30 |    |",
		"|    |  2 |    | 15 |",
		"",
	}

	res := out.String()
	assertRows(t, expectedLines, res)
}

func TestAddWithDynamicColumns(t *testing.T) {
	tab := New([]string{"a", "b", "c"})
	tab.Add(1, 2, 3)
	tab.Add(10, 20, 30)
	tab.Add("", 2)
	tab.Add("", "", "", "something")
	tab.Add("", "", "", 1.23)
	tab.Add("", "", "", func() string { return "cb-val" })
	out := bytes.NewBuffer(nil)
	err := tab.Print(out)
	if err != nil {
		t.Fatalf("Print failed: %s", err)
	}

	expectedLines := []string{
		"|  a |  b |  c |     col-3 |",
		"|----|----|----|-----------|",
		"|  1 |  2 |  3 |           |",
		"| 10 | 20 | 30 |           |",
		"|    |  2 |    |           |",
		"|    |    |    | something |",
		"|    |    |    |      1.23 |",
		"|    |    |    |    cb-val |",
		"",
	}

	res := out.String()
	assertRows(t, expectedLines, res)
}

func assertRows(t *testing.T, expectedLines []string, res string) {
	t.Helper()
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
	tab := New([]string{"a", "b", "c", "d"})
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

// ensure the stringerStruct implements fmt.Stringer at compile time
var _ fmt.Stringer = (*stringerStruct)(nil)

type stringerStruct struct {
	a string
	b bool
}

func (s *stringerStruct) String() string {
	return fmt.Sprintf("a=%q b=%t", s.a, s.b)
}

func TestValueString(t *testing.T) {

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name:  "nil",
			input: nil,
			want:  "",
		},
		{
			name:  "bool false",
			input: false,
			want:  "false",
		},
		{
			name:  "bool true",
			input: true,
			want:  "true",
		},
		{
			name:  "float32",
			input: float32(1.23),
			want:  "1.23",
		},
		{
			name:  "float64",
			input: float64(1.23),
			want:  "1.23",
		},
		{
			name:  "int",
			input: 123,
			want:  "123",
		},
		{
			name:  "int32",
			input: int32(123),
			want:  "123",
		},
		{
			name:  "int64",
			input: int64(123),
			want:  "123",
		},
		{
			name: "plain struct value",
			input: struct {
				a string
				b bool
			}{
				a: "aaa",
				b: true,
			},
			want: "{aaa true}",
		},
		{
			name: "plain struct pointer",
			input: &struct {
				a string
				b bool
			}{
				a: "aaa",
				b: true,
			},
			want: "&{aaa true}",
		},
		{
			name: "stringer struct value",
			input: &stringerStruct{
				a: "aaa",
				b: true,
			},
			want: `a="aaa" b=true`,
		},
		{
			name: "stringer struct pointer",
			input: &stringerStruct{
				a: "aaa",
				b: true,
			},
			want: `a="aaa" b=true`,
		},
		{
			name:  "callback no args",
			input: func() string { return "123" },
			want:  "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := valueString(tt.input)
			if res != tt.want {
				t.Errorf("valueString unexpected result: %q", res)
				t.Logf("diff (-want, +got): %s", cmp.Diff(tt.want, res))
			}
		})
	}
}

func ExampleTabulator_Add() {
	data := []struct {
		time   time.Time
		metric string
		count  any
	}{
		{
			time:   time.Time{},
			metric: "a",
			count:  1,
		},
		{
			time:   time.Time{},
			metric: "b",
			count:  20,
		},
		{
			time:   time.Time{},
			metric: "c",
			count:  1.23,
		},
	}

	tab := New([]string{"time", "metric", "count"})
	for _, elem := range data {
		tab.Add(func() string { return elem.time.Format("2006-01-02T15:04") }, elem.metric, elem.count)
	}

	tab.Print(os.Stdout)
	// Output:
	// |             time | metric | count |
	// |------------------|--------|-------|
	// | 0001-01-01T00:00 |      a |     1 |
	// | 0001-01-01T00:00 |      b |    20 |
	// | 0001-01-01T00:00 |      c |  1.23 |
}
