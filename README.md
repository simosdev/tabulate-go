# tabulate-go

## Installation

```shell
go get github.com/simosdev/tabulate-go
```

## Usage

```go
package main

import (
    "time"
    "os"

    "github.com/simosdev/tabulate-go"
)

func main() {
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
		tab.Add(
			func() string { return elem.time.Format("2006-01-02T15:04") },
			elem.metric,
			elem.count,
		)
	}

	tab.Print(os.Stdout)
	// Output:
	// |             time | metric | count |
	// |------------------|--------|-------|
	// | 0001-01-01T00:00 |      a |     1 |
	// | 0001-01-01T00:00 |      b |    20 |
	// | 0001-01-01T00:00 |      c |  1.23 |
}
```
