package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

	app = appStruct{
		Editor:     flag.String("editor", "subl", "Editor to use."),
		ConfigPath: flag.String("path", filepath.Dir(wd+"/")+"/", "runner.json path"),
		Port:       flag.String("port", "31777", "webserver port"),
		Gui:        flag.Bool("gui", true, "use gui window"),
	}
)

func main() {
	flag.Parse()

	config := &configStruct{Runners: []runnerStruct{}}
	if configPath, err := filepath.Abs(*app.ConfigPath + "/runner.json"); err == nil {
		if raw, err := ioutil.ReadFile(configPath); err == nil {
			if err = json.Unmarshal(raw, config); err != nil {
				fmt.Println("JSON error:", err)
			}
		} else {
			exit("runner.json not found")
		}
	}

	for id, service := range config.Runners {
		if !service.Ignore {
			makeService(id, *app.ConfigPath, service)
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
