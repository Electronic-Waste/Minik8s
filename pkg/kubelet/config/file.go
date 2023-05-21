package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"minik8s.io/pkg/constant"
	"os"
	"time"
)

type podEventType int

const (
	retryPeriod = 1 * time.Second
)

const (
	podAdd podEventType = iota
	podModify
	podDelete

	eventBufferLen = 10
)

type watchEvent struct {
	fileName  string
	eventType podEventType
}

type sourceFile struct {
	// using the channel in mux
	update chan interface{}

	// file channel
	watch chan *watchEvent

	// listening path
	path string
}

func NewSourceFile(ch chan interface{}) {
	cfg := newSourceFile(ch, constant.SysPodDir)
	cfg.run()
}

func newSourceFile(ch chan interface{}, path string) *sourceFile {
	return &sourceFile{
		update: ch,
		watch:  make(chan *watchEvent, eventBufferLen),
		path:   path,
	}
}

func (cfg *sourceFile) run() {
	go func() {
		// start to receive message from sourceFile
		select {
		case e := <-cfg.watch:
			{
				fmt.Println(e)
				cfg.update <- 1
			}
		}
	}()

	cfg.startWatch()
}

func (s *sourceFile) doWatch() error {
	_, err := os.Stat(s.path)
	if err != nil {
		fmt.Println("error in path of source")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("unable to create inotify: %v", err)
	}
	defer w.Close()

	err = w.Add(s.path)
	if err != nil {
		return fmt.Errorf("unable to create inotify for path %q: %v", s.path, err)
	}

	for {
		select {
		case event := <-w.Events:
			if err = s.produceWatchEvent(&event); err != nil {
				return fmt.Errorf("error while processing inotify event (%+v): %v", event, err)
			}
		case err = <-w.Errors:
			return fmt.Errorf("error while watching %q: %v", s.path, err)
		}
	}
}

func (s *sourceFile) produceWatchEvent(e *fsnotify.Event) error {
	var eventType podEventType
	switch {
	case (e.Op & fsnotify.Create) > 0:
		eventType = podAdd
	case (e.Op & fsnotify.Write) > 0:
		eventType = podModify
	case (e.Op & fsnotify.Chmod) > 0:
		eventType = podModify
	case (e.Op & fsnotify.Remove) > 0:
		eventType = podDelete
	case (e.Op & fsnotify.Rename) > 0:
		eventType = podDelete
	default:
		// Ignore rest events
		return nil
	}

	s.watch <- &watchEvent{e.Name, eventType}
	return nil
}

func (cfg *sourceFile) startWatch() {
	go func() {
		for {
			cfg.doWatch()
			time.Sleep(retryPeriod)
		}
	}()
}
