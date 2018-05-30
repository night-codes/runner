package main

import (
	"github.com/night-codes/tokay"
	"github.com/night-codes/webview"
	"path/filepath"
	"runtime"
	"strconv"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func webserver(config *configStruct) {
	r := tokay.New(&tokay.Config{TemplatesDirs: []string{basepath + "/templates"}})
	r.Debug = false

	r.Static("/files", basepath+"/files")

	r.GET("/", func(c *tokay.Context) {
		c.Redirect(303, "/service/0")
	})

	r.GET("/service/<active:\\d+>", func(c *tokay.Context) {
		c.HTML(200, "index", map[string]interface{}{
			"services": activeServices,
			"ignored":  ignoredServices,
			"config":   config,
			"active":   c.ParamInt("active"),
		})
	})

	r.GET("/logs/<active:\\d+>", func(c *tokay.Context) {
		active := c.ParamInt("active")
		for _, service := range activeServices {
			if service.ID == active {
				if _, ex := c.QueryEx("json"); !ex {
					logs := ""
					for _, log := range service.Logs {
						logs += log.Time.Format("[2006-01-02 15:04:05] ") + log.Message + "\n"
					}
					c.String(200, logs)
				} else {
					c.JSON(200, service.Logs)
				}
				return
			}
		}
		c.String(404, "Not found")
	})

	// GUI start
	port := strconv.Itoa(app.Port)
	if app.Gui {
		go r.Run(":" + port)
		w := webview.New(webview.Settings{
			Title:     "Runner",
			Icon:      basepath + "/files/img/favicon.png",
			URL:       "http://localhost:" + port,
			Height:    800,
			Width:     1200,
			Resizable: true,
		})
		w.SetColor(73, 82, 88, 255)
		w.Run()
	} else {
		r.Run(":" + port)
	}
}
