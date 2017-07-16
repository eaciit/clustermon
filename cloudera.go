package clustermon

type Cloudera struct {
	Cluster
}

func NewCloudera(name string, uri string) ICluster {
	c := new(Cloudera)
	c.Name = name
	c.APIUri = uri
	return c
}
