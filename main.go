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
		Title   string         `json:"title"`
		Runners []runnerStruct `json:"services"`
	}
)

func main() {
	flag.Parse()

	config := &configStruct{Runners: []runnerStruct{}}
	if configPath, err := filepath.Abs(*cpath + "/runner.json"); err == nil {
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
			makeService(id, *cpath, service)
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
