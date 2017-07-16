package clustermon

type Service struct{
	Name string
	Metrics map[string]*Metric
}