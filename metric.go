package clustermon

import "time"

type Metric struct {
	Name string
}

type TimeValue struct {
	Time   time.Time
	Values map[string]float64
}
