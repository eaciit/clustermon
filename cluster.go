package clustermon

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"encoding/json"
	"net/http"

	"github.com/eaciit/toolkit"
)

type ICluster interface {
	ClusterName() string
	HostNames() []string
	ServiceNames() []string
	MetricNames(string) []string
	Values(hostname, servicename string, fields map[string]string, from, to time.Time, interval int) ([]toolkit.M, error)
	RefreshMeta() error
}

type Cluster struct {
	Name   string
	Active bool

	UserName, Password string
	Hosts              []*Host
	Services           []*Service

	Provider string
	APIUri   string
}

func (c *Cluster) ClusterName() string {
	return c.Name
}

func (c *Cluster) HostNames() []string {
	ret := []string{}
	for _, v := range c.Hosts {
		ret = append(ret, v.Name)
	}
	return ret
}
func (c *Cluster) ServiceNames() []string {
	ret := []string{}
	for _, v := range c.Services {
		ret = append(ret, v.Name)
	}
	return ret
}
func (c *Cluster) MetricNames(string) []string { return []string{} }
func (c *Cluster) RefreshMeta() error          { return fmt.Errorf("RefreshMete is not yet implemented") }

func (c *Cluster) Values(hostname, servicename string, fields map[string]string, from, to time.Time, interval int) ([]toolkit.M, error) {
	var res []toolkit.M
	return res, fmt.Errorf("Not yet implemented")
}

func (c *Cluster) call(apiPath string, data []byte) (map[string]interface{}, error) {
	uri := c.APIUri
	m := make(map[string]interface{})
	if apiPath != "" {
		if !strings.HasSuffix(uri, "/") {
			uri += "/"
		}
		if strings.HasPrefix(apiPath, "/") && len(apiPath) >= 2 {
			apiPath = apiPath[1:]
		}
		uri += apiPath

		r, err := toolkit.HttpCall(uri, http.MethodGet, data,
			toolkit.M{}.Set("auth", "basic").Set("user", c.UserName).Set("password", c.Password))
		if err != nil {
			return m, fmt.Errorf("Unable to call %s: %s", uri, err.Error())
		}
		defer r.Body.Close()

		output, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return m, fmt.Errorf("Unable to call %s: %s", uri, err.Error())
		}

		if err = json.Unmarshal(output, &m); err != nil {
			return m, fmt.Errorf("Unable to call %s: %s", uri, err.Error())
		}
		return m, nil
	}
	return m, nil
}
