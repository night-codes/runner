package main

import (
	"log"
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
)

func webserver(config *configStruct) {
	r := tokay.New(&tokay.Config{TemplatesDirs: []string{basepath + "/templates"}})
	r.Debug = true

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

	wss := ws.NewTokay("/ws/connect<num:\\d+>", &r.RouterGroup)

	wss.Read("test", func(a *ws.Adapter) {
		log.Println(string(a.Command()))
		a.Send("message interface{}")
	})
	wss.Read("test2", func(a *ws.Adapter) {
		log.Println(string(a.Command()))
		time.Sleep(time.Second * 10)
		log.Println(a.Connection().Request("p4", map[string]interface{}{"task": "build", "status": "OK", "time": 10.45}, 10))
	})
	wss.Read("close", func(a *ws.Adapter) {
		a.Send("OK")
		a.Close()
	})
	go func() {
		for t := range time.Tick(time.Second * 3) {
			wss.Send("ololo", t)
		}
	}()

	go func() {
		for t := range time.Tick(time.Second * 3) {
			wss.Subscribers("news").Send("news", ws.Map{"time": t, "message": "news"})
		}
	}()

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
