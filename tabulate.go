package tabulate

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Row map[string]any

type tabulater struct {
	cols []string
	rows []Row
}

func New(cols ...string) *tabulater {
	res := &tabulater{}
	res.Columns(cols...)
	return res
}

func (t *tabulater) Columns(cols ...string) {
	t.cols = append(t.cols, cols...)
}

func (t *tabulater) AddRow(row Row) {
	// TODO: add possible new keys to cols slice if not already present
	t.rows = append(t.rows, row)
}

func (t *tabulater) Print(w io.Writer) error {
	maxValueLengths := t.colMaxValueLengths()

	// column names
	for _, col := range t.cols {
		diff := maxValueLengths[col] - len(col)
		diff = max(diff, 0)
		_, err := fmt.Fprintf(w, "| %s%s ", strings.Repeat(" ", diff), col)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "|\n")
	if err != nil {
		return err
	}

	// columns second row
	for _, col := range t.cols {
		_, err := fmt.Fprintf(w, "|%s", strings.Repeat("-", maxValueLengths[col]+2))
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, "|\n")
	if err != nil {
		return err
	}

	for _, row := range t.rows {
		for _, col := range t.cols {
			diff := maxValueLengths[col] - valueLength(row[col])
			diff = max(diff, 0)
			_, err := fmt.Fprintf(w, "| %s%s ", strings.Repeat(" ", diff), valueString(row[col]))
			if err != nil {
				return err
			}
		}
		_, err := fmt.Fprintf(w, "|\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *tabulater) colMaxValueLengths() map[string]int {
	lengths := make(map[string]int, 0)
	for _, row := range t.rows {
		for _, col := range t.cols {
			if oldLen, ok := lengths[col]; ok {
				vl := valueLength(row[col])
				if vl > oldLen {
					lengths[col] = vl
				}
			} else {
				lengths[col] = valueLength(row[col])
			}
		}
	}
	return lengths
}

func valueLength(val any) int {
	return len(valueString(val))
}

func valueString(val any) string {
	if val == nil {
		return ""
	}

	var res string
	switch val := val.(type) {
	case float32:
		res = strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		res = strconv.FormatFloat(val, 'f', -1, 64)
	case fmt.Stringer:
		res = val.String()
	default:
		res = fmt.Sprintf("%v", val)
	}
	return res
}
