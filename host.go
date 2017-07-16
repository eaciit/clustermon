package clustermon

type Host struct {
	Name     string
	Services map[string]*Service
}
