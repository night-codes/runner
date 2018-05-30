package main

import (
	"encoding/json"
	"fmt"
	"github.com/alexflint/go-arg"
	"io/ioutil"
	"os"
	"path/filepath"
)

type (
	configStruct struct {
		Title   string         `json:"title"`
		Runners []runnerStruct `json:"services"`
	}
	appStruct struct {
		Editor     *string
		ConfigPath *string
		Port       *string
		Gui        *bool
	}
)

var (
	wd, _ = os.Getwd()

	app struct {
		ConfigPath string `arg:"positional" help:"runner.json path"`
		Editor     string `arg:"-e" help:"editor to use"`
		Port       int    `arg:"-p" help:"webserver port"`
		Gui        bool   `arg:"-g" help:"use gui window"`
	}
)

func main() {
	app.ConfigPath = dir(wd)
	app.Editor = "subl"
	app.Port = 31777
	app.Gui = true

	arg.MustParse(&app)

	config := &configStruct{Runners: []runnerStruct{}}
	configPath := configpath(app.ConfigPath)

	if raw, err := ioutil.ReadFile(configPath); err == nil {
		if err = json.Unmarshal(raw, config); err != nil {
			exit("JSON error: " + err.Error())
		}
	} else {
		exit("Can't read runner.json")
	}

	for id, service := range config.Runners {
		if !service.Ignore {
			makeService(id, filepath.Dir(configPath), service)
		} else {
			ignoredServices = append(ignoredServices, &serviceStruct{runnerStruct: service, ID: id})
		}
	}

	webserver(config)
}

func exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}

func dir(dir string) string {
	return filepath.Dir(dir+"/") + "/"
}

func configpath(path string) string {
	path, _ = filepath.Abs(path)
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return path
		}
		if path, err = filepath.Abs(path + "/runner.json"); err == nil {
			return path
		}
	}
	exit("runner.json not found")
	return ""
}
