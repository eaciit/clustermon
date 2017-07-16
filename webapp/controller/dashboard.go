package controller

import (
	"eaciit/clustermon"
	"time"

	"fmt"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

var (
	clusters []clustermon.ICluster
)

func SetCluster(cs []clustermon.ICluster) {
	clusters = cs
}

type Dashboard struct {
}

func (d *Dashboard) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate

	/*
		jsonClusters := []toolkit.M{}
		for _, c := range clusters {
			jsonc := toolkit.M{}
			jsonc.Set("name", c.ClusterName())
			jsonClusters = append(jsonClusters, jsonc)
		}
	*/

	return toolkit.M{}.Set("cluster", clusters)
}

func (d *Dashboard) GetClusters(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	return clusters
}

func (d *Dashboard) GetCPU(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	result := toolkit.NewResult()
	clusterid := ctx.Query("name")
	c, err := findCluster(clusterid)
	if err != nil {
		return result.SetError(err)
	}

	t0 := time.Now().Add(-1 * time.Hour)
	t1 := time.Now()
	values, err := c.Values("", "", map[string]string{
		"metrics/cpu/System._avg": "System",
		"metrics/cpu/User._avg":   "User",
	}, t0, t1, 15)
	if err != nil {
		return result.SetError(err)
	}

	result.SetData(values)
	return result
}

func (d *Dashboard) GetMemory(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	result := toolkit.NewResult()
	clusterid := ctx.Query("name")
	c, err := findCluster(clusterid)
	if err != nil {
		return result.SetError(err)
	}

	t0 := time.Now().Add(-1 * time.Hour)
	t1 := time.Now()
	values, err := c.Values("", "", map[string]string{
		"metrics/memory/Share._avg":  "Share",
		"metrics/memory/Use._avg":    "Use",
		"metrics/memory/Cache._avg":  "Cache",
		"metrics/memory/Swap._avg":   "Swap",
		"metrics/memory/Buffer._avg": "Buffer",
	}, t0, t1, 15)
	if err != nil {
		return result.SetError(err)
	}

	result.SetData(values)
	return result
}

func (d *Dashboard) GetNetwork(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	result := toolkit.NewResult()
	clusterid := ctx.Query("name")
	c, err := findCluster(clusterid)
	if err != nil {
		return result.SetError(err)
	}

	t0 := time.Now().Add(-1 * time.Hour)
	t1 := time.Now()
	values, err := c.Values("", "", map[string]string{
		"metrics/network/Out._avg": "Out",
		"metrics/network/In._avg":  "In",
	}, t0, t1, 15)
	if err != nil {
		return result.SetError(err)
	}

	result.SetData(values)
	return result
}

func findCluster(name string) (clustermon.ICluster, error) {
	for _, c := range clusters {
		if c.ClusterName() == name {
			return c, nil
		}
	}

	return nil, fmt.Errorf("Cluster %s could not be found", name)
}
