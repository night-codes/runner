package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"gopkg.in/night-codes/types.v1"

	"github.com/fatih/color"
	"github.com/night-codes/ws"
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

func init() {
	mainWS.Read("stop", func(a *ws.Adapter) {
		activeServices[types.Int(a.Data())].stop()
	})

	mainWS.Read("clear", func(a *ws.Adapter) {
		activeServices[types.Int(a.Data())].Logs = []logMessage{}
	})

	mainWS.Read("start", func(a *ws.Adapter) {
		activeServices[types.Int(a.Data())].start()
	})

	mainWS.Read("restart", func(a *ws.Adapter) {
		activeServices[types.Int(a.Data())].stop()
		time.Sleep(time.Second / 3)
		activeServices[types.Int(a.Data())].start()
	})
}

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
	s.makeCmd()
	if !s.Stopped {
		go func(s *serviceStruct) {
			time.Sleep(time.Second / 2)
			if s.Delay > 0 {
				fmt.Printf("Delay %v ms before %v starting...\n", s.Delay, s.Title)
				time.Sleep(time.Millisecond * time.Duration(s.Delay))
			}
			s.runService()
		}(s)
	}
	return s
}

func (s *serviceStruct) runService() {
	s.changeStatus(statusWaiting, false)
	green := color.New(color.FgGreen, color.Bold)
	green.Printf("Started: %v\n", s.Title)
	err := s.Cmd.Run()

	red := color.New(color.FgRed, color.Bold)
	if err != nil {
		red.Printf("%v finished with: %v\n", s.Title, err)
		s.ErrLogger.WriteString(err.Error())
	} else {
		red.Printf("Finished: %v\n", s.Title)
	}
	s.changeStatus(statusStopped)
	s.TitleLogger.WriteString(s.Title + " finished\n")

	if !s.Stopped && s.Restart {
		if s.RestartDelay > 0 {
			fmt.Printf("Delay %v ms before %v restarting...\n", s.RestartDelay, s.Title)
			time.Sleep(time.Millisecond * time.Duration(s.RestartDelay))
		}
		s.makeCmd()
		s.runService()
	}
}

func (s *serviceStruct) changeStatus(status int, stopped ...bool) {
	s.Status = status
	if len(stopped) > 0 {
		s.Stopped = stopped[0]
	}
	mainWS.Send("changeStatus", obj{"status": status, "service": s.ID})
}

func (s *serviceStruct) makeCmd() {
	s.Cmd = exec.Command("bash", "-c", s.Command)
	s.Cmd.Dir = s.Dir
	s.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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

func (s *serviceStruct) stop() {
	if s.Cmd != nil && (s.Status == statusRunned || s.Status == statusWaiting) {
		pgid, err := syscall.Getpgid(s.Cmd.Process.Pid)
		if err == nil {
			s.changeStatus(statusWaiting, true)
			syscall.Kill(-pgid, syscall.SIGKILL)
		}
		s.Cmd.Wait()
	}
}

func (s *serviceStruct) start() {
	if s.Cmd != nil && s.Status == statusStopped {
		s.makeCmd()
		s.runService()
	}
}

/**
 * Opening file:14:3 in different editors:
 *
 * sublime_text ./file:14:3
 * atom ./file:14:3
 * code -g ./file:14:3
 *
 */
