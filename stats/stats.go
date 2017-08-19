package stats

import (
    "fmt"
    "time"
)

type Node struct {
    Uuid    string
    Stats   *Status
}

func (nd Node) String() string {
    return "todo String()"
}

type ThroughPut struct {
    Count   int64
    Len     int64
    AvgLat  time.Duration
    MinLat  time.Duration
    MaxLat  time.Duration
}

func (tp ThroughPut) String() string {
    return fmt.Sprintf("\tCount: %v\n\tLen: %v\n\tAvgLat: %v\n", tp.Count, tp.Len, tp.AvgLat)
}

type Status struct {
    X200 *ThroughPut
    X400 *ThroughPut
    X500 *ThroughPut
}

func (s Status) String() string {
    return fmt.Sprintf("2XX:\n%v\n4XX:\n%v\n5XX\n%v\n", s.X200, s.X400, s.X500)
}
