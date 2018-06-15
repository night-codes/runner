package main

import (
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/night-codes/tokay"
	"github.com/night-codes/webview"
	"github.com/night-codes/ws"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	r          = tokay.New(&tokay.Config{TemplatesDirs: []string{basepath + "/templates"}})
	mainWS     = ws.NewTokay("/ws", &r.RouterGroup)
	wv         webview.WebView
)

func webserver(config *configStruct) {
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
				c.JSON(200, service.Logs)
				return
			}
		}
		c.String(404, "Not found")
	})

	// GUI start
	port := strconv.Itoa(app.Port)
	if app.Gui {
		go r.Run(":" + port)
		wv = webview.New(webview.Settings{
			Title:     "Runner",
			Icon:      basepath + "/files/img/favicon.png",
			URL:       "http://localhost:" + port,
			Height:    800,
			Width:     1200,
			Resizable: true,
		})
		wv.SetColor(73, 82, 88, 255)
		wv.Run()
	} else {
		r.Run(":" + port)
	}
	for _, s := range activeServices {
		s.stop()
	}
	time.Sleep(time.Second / 2)
}
