package clustermon

import (
	"sort"
	"strings"
	"time"

	"github.com/eaciit/toolkit"
)
import "fmt"

type AmbariV1 struct {
	Cluster
}

func NewAmbariV1(name, uri, user, password string) ICluster {
	c := new(AmbariV1)
	c.Name = name
	c.APIUri = uri
	c.UserName = user
	c.Password = password
	return c
}

func (c *AmbariV1) Metric(hostname, servicename, metricname string, from, to time.Time) ([]*Metric, error) {
	var res []*Metric
	return res, nil
}

func (c *AmbariV1) RefreshMeta() error {
	var (
		output map[string]interface{}
		err    error
	)

	//--- get host
	if output, err = c.call(fmt.Sprintf("/clusters/%s/hosts", c.Name), nil); err != nil {
		return fmt.Errorf("Get hosts failed: %s", err.Error())
	}
	items := output["items"].([]interface{})
	c.Hosts = []*Host{}
	for _, item := range items {
		i := item.(map[string]interface{})["Hosts"].(map[string]interface{})
		h := new(Host)
		h.Name = i["host_name"].(string)
		c.Hosts = append(c.Hosts, h)
	}

	//--- get services
	if output, err = c.call(fmt.Sprintf("/clusters/%s/services", c.Name), nil); err != nil {
		return fmt.Errorf("Get hosts failed: %s", err.Error())
	}
	items = output["items"].([]interface{})
	c.Services = []*Service{}
	for _, item := range items {
		i := item.(map[string]interface{})["ServiceInfo"].(map[string]interface{})
		s := new(Service)
		s.Name = i["service_name"].(string)
		c.Services = append(c.Services, s)
	}

	return nil
}

func (c *AmbariV1) Values(hostname, servicename string, fields map[string]string, from, to time.Time, interval int) ([]toolkit.M, error) {
	var res []toolkit.M
	fieldparms := []string{}
	for kfield, _ := range fields {
		fieldparm := fmt.Sprintf("%s[%d,%d,%d]", kfield, from.Unix(), to.Unix(), interval)
		fieldparms = append(fieldparms, fieldparm)
	}

	apiPath := "/clusters/" + c.Name
	if hostname != "" {
		apiPath += "/hosts/" + hostname
	}
	if servicename != "" {
		apiPath += "/services/" + servicename
	}
	fieldqueries := strings.Join(fieldparms, ",")
	apiPath += "?fields=" + fieldqueries

	output, err := c.call(apiPath, nil)
	if err != nil {
		return res, err
	}

	datamap := map[time.Time]toolkit.M{}
	metrics := output["metrics"].(map[string]interface{})
	for ktag, ttag := range fields {
		if ttag == "" {
			ttag = ktag
		}
		ktag = strings.TrimPrefix(ktag, "metrics/")
		ms := getMetricesValues(metrics, ktag, ttag)
		for ts, val := range ms {
			m, ok := datamap[ts]
			if !ok {
				m = toolkit.M{}
				m.Set("Time", ts)
				datamap[ts] = m
			}

			m.Set(ttag, val)
		}
	}

	for _, v := range datamap {
		res = append(res, v)
	}

	sort.Sort(ValueSorter(res))
	return res, nil
}

func getMetricesValues(tablemap map[string]interface{}, tag string, title string) map[time.Time]float64 {
	res := map[time.Time]float64{}
	tags := strings.Split(tag, "/")
	firsttag := tags[0]
	if len(tags) > 1 {
		nexttag := strings.Join(tags[1:], "/")
		return getMetricesValues(tablemap[firsttag].(map[string]interface{}), nexttag, title)
	} else {
		values := tablemap[firsttag].([]interface{})
		for _, v0 := range values {
			v1 := v0.([]interface{})
			value, timestamp := v1[0].(float64), time.Unix(int64(v1[1].(float64)), 0)
			res[timestamp] = value
		}
	}
	return res
}

type ValueSorter []toolkit.M

func (vs ValueSorter) Len() int {
	return len(vs)
}

func (vs ValueSorter) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs ValueSorter) Less(i, j int) bool {
	return vs[i].Get("Time").(time.Time).Before(vs[j].Get("Time").(time.Time))
}
