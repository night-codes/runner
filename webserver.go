package main

import (
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/night-codes/tokay"
	"github.com/night-codes/tokay-ws"
	"github.com/night-codes/webview"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func webserver(config *configStruct) {
	r := tokay.New(&tokay.Config{TemplatesDirs: []string{basepath + "/templates"}})
	r.Debug = false
	ws1 := ws.New("/ws/connect", &r.RouterGroup)
	ws2 := ws.New("/ws/connect2", &r.RouterGroup)

	ws1.Read("test", func(a *ws.Adapter) {
		log.Println(string(a.Command()))
		a.Send("message interface{}")
	})
	ws2.Read("test2", func(a *ws.Adapter) {
		log.Println(string(a.Command()))
	})
	ws2.Read("close", func(a *ws.Adapter) {
		a.Send("OK")
		a.Close()
	})

	go func() {
		for t := range time.Tick(time.Second * 3) {
			ws1.Send("ololo", t)
		}
	}()

	go func() {
		for t := range time.Tick(time.Second * 3) {
			ws2.Subscribers("news").Send("news", ws.Map{"time": t, "message": "news"})
		}
	}()

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
