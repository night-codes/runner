package main

import (
	"github.com/night-codes/tokay"
	"github.com/night-codes/webview"
	"path/filepath"
	"runtime"
	"fmt"
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
		for _,service := range activeServices {
			fmt.Println(service.ID, active)
			if service.ID == active {
				c.JSON(200, service.Logs)
				return
			}
		}
		c.String(404, "Not found")
	})





	// GUI start
	if *app.Gui {
		go r.Run(":" + *app.Port)
		w := webview.New(webview.Settings{
			Title:     "Runner",
			Icon:      basepath + "/files/img/favicon.png",
			URL:       "http://localhost:" + *app.Port,
			Height:    800,
			Width:     1200,
			Resizable: true,
		})
		w.SetColor(73, 82, 88, 255)
		w.Run()
	} else {
		r.Run(":" + *app.Port)
	}
}
