package config

import (
	"fmt"
	"minik8s.io/pkg/constant"
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
			}
		}
	}()

	cfg.startWatch()
}

func (cfg *sourceFile) startWatch() {
	go func() {
		// do watch

		// sleep for a time
		time.Sleep(retryPeriod)
	}()
}
