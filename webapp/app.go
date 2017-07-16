package main

import (
	"eaciit/clustermon/webapp/controller"
	"os"
	"path/filepath"

	"github.com/eaciit/config"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"

	"eaciit/clustermon"
)

var (
	exepath  = ""
	monitors = []clustermon.ICluster{}
)

func ExePath() string {
	if exepath == "" {
		exePath, _ := os.Executable()
		exePath = filepath.Dir(exePath)
	}
	return exepath
}

func main() {
	log, _ := toolkit.NewLog(true, false, "", "", "")
	configpath := filepath.Join(ExePath(), "..", "config", "app.json")
	econfig := config.SetConfigFile(configpath)
	if econfig != nil {
		log.Error("Error loading config file " + econfig.Error())
	}

	//--- clusters
	clusters := config.Get("clusters").([]interface{})
	for _, cface := range clusters {
		var cluster clustermon.Cluster
		toolkit.Serde(cface, &cluster, "")
		if cluster.Active {
			if cluster.Provider == "AmbariV1" {
				ambariv1 := &clustermon.AmbariV1{Cluster: cluster}
				monitors = append(monitors, ambariv1)
			} else if cluster.Provider == "Cloudera" {
				cloudera := &clustermon.Cloudera{Cluster: cluster}
				monitors = append(monitors, cloudera)
			}
			log.Infof("Monitor cluster %s:%s:%s", cluster.Provider, cluster.Name, cluster.APIUri)
		}
	}
	for _, m := range monitors {
		m.RefreshMeta()
	}
	controller.SetCluster(monitors)

	port := int(config.GetDefault("port", 9100).(float64))
	serveraddress := config.GetDefault("server", "0.0.0.0").(string)

	wd := config.GetDefault("workingpath", "").(string)
	app := App(wd)
	knot.StartApp(app, toolkit.Sprintf("%s:%d", serveraddress, port))
	/*
		knot.StartAppWithFn(app, toolkit.Sprintf("%s:%d", serveraddress, port),
			map[string]knot.FnContent{
				"/": func(r *knot.WebContext) interface{} {
					http.Redirect(r.Writer, r.Request, "/dashboard/index", 301)
					return nil
				}})
	*/
}

func App(wd string) *knot.App {
	app := knot.NewApp("BigData Cluster Monitor")
	if wd == "" {
		wd = filepath.Join(ExePath(), "..", "webapp")
	}
	app.ViewsPath = filepath.Join(wd, "views")
	app.LayoutTemplate = "_layout.html"
	app.Static("static", filepath.Join(wd, "assets"))
	app.Static("views", filepath.Join(wd, "views"))

	// Register the app
	app.Register(new(controller.Dashboard))

	app.DefaultOutputType = knot.OutputHtml
	return app
}
