package main

import (
	"log"
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
		Title        string   `json:"title"`
		Path         string   `json:"path"`
		After        []string `json:"after"`
		StartMessage string   `json:"startMessage"`
		Restart      bool     `json:"restart"`
		Stopped      bool     `json:"stopped"`
		Ignore       bool     `json:"ignore"`
		Delay        uint64   `json:"delay"`
		RestartDelay uint64   `json:"restartDelay"`
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
	s.TitleLogger = logger{Type: typeInfo, Service: s}
	s.InfoLogger = logger{Type: typeTitle, Service: s}
	s.ErrLogger = logger{Type: typeError, Service: s}
	s.TitleLogger.WriteString("# " + s.Dir + "\n")
	s.Dir, _ = filepath.Abs(basepath + "/" + s.Path + "/")
	s.ID = id
	activeServices = append(activeServices, s)
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
		log.Printf("%v finished with: %v", s.Title, err)
		s.ErrLogger.WriteString(err.Error())
	} else {
		log.Printf("Finished: %v", s.Title)
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
	s.Cmd = exec.Command("bash", "-c", "go run *.go")
	s.Cmd.Dir = s.Dir
	s.Cmd.Stdout = &s.InfoLogger
	s.Cmd.Stderr = &s.ErrLogger
}

/**
 * Opening file:14:3 in different editors:
 *
 * sublime_text ./file:14:3
 * atom ./file:14:3
 * code -g ./file:14:3
 *
 */
