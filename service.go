package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	statusStopped = 0
	statusWaiting = 1
	statusRunned  = 2
	statusIgnored = 3
)

type (
	logMessage struct {
		Time    time.Time
		Message string
		Type    int
	}
	runnerStruct struct {
		Title        string            `json:"title"`
		Path         string            `json:"path"`
		After        []string          `json:"after"`
		StartMessage string            `json:"startMessage"`
		Restart      bool              `json:"restart"`
		Stopped      bool              `json:"stopped"`
		Ignore       bool              `json:"ignore"`
		Delay        uint64            `json:"delay"`
		RestartDelay uint64            `json:"restartDelay"`
		Command      string            `json:"command"`
		Env          map[string]string `json:"env"`
	}
	serviceStruct struct {
		runnerStruct
		ID          int
		Cmd         *exec.Cmd
		Logs        []logMessage
		TitleLogger logger
		InfoLogger  logger
		ErrLogger   logger
		Status      int
		Dir         string
	}
)

var (
	activeServices  = []*serviceStruct{}
	ignoredServices = []*serviceStruct{}
)

func makeService(id int, basepath string, r runnerStruct) *serviceStruct {
	s := &serviceStruct{runnerStruct: r}
	s.Logs = []logMessage{}
	s.TitleLogger = logger{Type: typeTitle, Service: s}
	s.InfoLogger = logger{Type: typeInfo, Service: s}
	s.ErrLogger = logger{Type: typeError, Service: s}
	s.Dir, _ = filepath.Abs(basepath + "/" + s.Path + "/")
	s.TitleLogger.WriteString("# " + s.Dir + "\n")
	s.ID = id
	activeServices = append(activeServices, s)
	if s.Command == "" {
		s.Command = "go build -o build;./build;rm build"
	}
	makeCmd(s)
	if !s.Stopped {
		go func(s *serviceStruct) {
			if s.Delay > 0 {
				log.Printf("Delay %v ms before %v starting...", s.Delay, s.Title)
				time.Sleep(time.Millisecond * time.Duration(s.Delay))
			}
			runService(s)
		}(s)
	}
	return s
}

func runService(s *serviceStruct) {
	log.Printf("Started: %v", s.Title)
	s.Stopped = false
	s.Status = statusWaiting
	err := s.Cmd.Run()
	if err != nil {
		log.Printf("%v finished with: %v\n", s.Title, err)
		s.ErrLogger.WriteString(err.Error())
	} else {
		log.Printf("Finished: %v\n", s.Title)
	}
	s.Status = statusStopped
	s.TitleLogger.WriteString(s.Title + " finished\n")

	if !s.Stopped && s.Restart {
		if s.RestartDelay > 0 {
			log.Printf("Delay %v ms before %v restarting...", s.RestartDelay, s.Title)
			time.Sleep(time.Millisecond * time.Duration(s.RestartDelay))
		}
		makeCmd(s)
		runService(s)
	}
}

func makeCmd(s *serviceStruct) {
	s.Cmd = exec.Command("bash", "-c", s.Command)
	s.Cmd.Dir = s.Dir
	s.Cmd.Stdout = &s.InfoLogger
	s.Cmd.Stderr = &s.ErrLogger
	e := []string{}
	for _, v := range os.Environ() {
		e = append(e, v)
	}
	for k, v := range s.Env {
		e = append(e, k+"="+v)
	}
	s.Cmd.Env = e
}

/**
 * Opening file:14:3 in different editors:
 *
 * sublime_text ./file:14:3
 * atom ./file:14:3
 * code -g ./file:14:3
 *
 */
