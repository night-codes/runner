package main

import (
	"github.com/night-codes/tokay"
	"github.com/night-codes/webview"
	"path/filepath"
	"runtime"
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
	go r.Run(":" + *port)
	w := webview.New(webview.Settings{
		Title:     "Runner",
		Icon:      basepath + "/files/img/favicon.png",
		URL:       "http://localhost:" + *port,
		Height:    800,
		Width:     1200,
		Resizable: true,
	})
	w.SetColor(73, 82, 88, 255)
	/* 	w.Dispatch(func() {
		w.Dialog(webview.DialogTypeAlert, webview.DialogFlagInfo, "test", "test")
		log.Println("Dispath!")
	}) */

	w.Run()
}
