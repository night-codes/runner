package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	editor = flag.String("editor", "sublime", "Editor to use.")
	ex, _  = os.Getwd()
	cpath  = flag.String("path", filepath.Dir(ex+"/")+"/", "runner.json path")
	port   = flag.String("port", "31777", "webserver port")
)

type (
	configStruct struct {
		Title   string                  `json:"title"`
		Runners map[string]runnerStruct `json:"services"`
	}
)

func main() {
	flag.Parse()

	config := &configStruct{Runners: map[string]runnerStruct{}}
	if configPath, err := filepath.Abs(*cpath + "/runner.json"); err == nil {
		if raw, err := ioutil.ReadFile(configPath); err == nil {
			if err = json.Unmarshal(raw, config); err != nil {
				fmt.Println("JSON error:", err)
			}
		} else {
			exit("runner.json not found")
		}
	}

	active := ""
	for name, service := range config.Runners {
		if !service.Ignore {
			if active == "" {
				active = name
			}
			makeService(name, *cpath, service)
		} else {
			ignoredServices = append(ignoredServices, &serviceStruct{runnerStruct: service, Name: name})
		}
	}

	if active == "" {
		exit("No active services specified in runner.json")
	}
	webserver(config, active)
}

func exit(err string) {
	fmt.Println(err)
	os.Exit(1)
}
