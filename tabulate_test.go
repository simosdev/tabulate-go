package tabulate

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateTableWithInstance(t *testing.T) {
	tab := New([]string{"a", "bbb", "c", "d"})
	tab.AddRow(Row{"a": 1, "bbb": 2, "c": 3})
	tab.AddRow(Row{"a": 10, "bbb": 20, "c": 30})
	tab.AddRow(Row{"bbb": 2})
	tab.AddRow(Row{"d": "something"})
	tab.AddRow(Row{"d": 1.23})
	tab.AddRow(Row{"d": func() string { return "cb-val" }})
	out := bytes.NewBuffer(nil)
	err := tab.Print(out)
	if err != nil {
		t.Fatalf("Print failed: %s", err)
	}

	expectedLines := []string{
		"|  a | bbb |  c |         d |",
		"|----|-----|----|-----------|",
		"|  1 |   2 |  3 |           |",
		"| 10 |  20 | 30 |           |",
		"|    |   2 |    |           |",
		"|    |     |    | something |",
		"|    |     |    |      1.23 |",
		"|    |     |    |    cb-val |",
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
