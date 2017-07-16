package test

import (
	"eaciit/clustermon"
	"fmt"
	"testing"
	"time"

	"github.com/eaciit/toolkit"
)

var hdp *clustermon.AmbariV1

func TestConnect(t *testing.T) {
	hdp = clustermon.NewAmbariV1("echdp01", "http://35.186.145.74:8080/api/v1", "admin", "Hulk.1234").(*clustermon.AmbariV1)
	if err := hdp.RefreshMeta(); err != nil {
		t.Errorf("Unable to refresh meta: %s", err.Error())
	} else {
		fmt.Printf("Cluster has %d hosts: %s\n", len(hdp.Hosts), hdp.HostNames())
		fmt.Printf("Cluster has %d services: %s\n", len(hdp.Services), hdp.ServiceNames())
	}
}

func TestGetValue(t *testing.T) {
	t0 := time.Now().Add(-3 * time.Hour)
	t1 := time.Now()
	values, err := hdp.Values("", "", map[string]string{
		"metrics/cpu/Nice._avg":   "CPU Nice(Avg)",
		"metrics/cpu/Idle._avg":   "CPU Idle(Avg)",
		"metrics/cpu/User._avg":   "CPU User(Avg)",
		"metrics/cpu/System._avg": "CPU System(Avg)",
	}, t0, t1, 15)
	if err != nil {
		t.Errorf("Fail: %s", err.Error())
	} else {
		fmt.Printf("Result: %s\n", toolkit.JsonStringIndent(values, "\t"))
	}
}
