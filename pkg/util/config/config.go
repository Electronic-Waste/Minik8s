package config

import (
	"fmt"
	"sync"
)

type Merger interface {
	// Invoked when a change from a source is received.  May also function as an incremental
	// merger if you wish to consume changes incrementally.  Must be reentrant when more than
	// one source is defined.
	Merge(source string, update interface{}) error
}

// MergeFunc implements the Merger interface
type MergeFunc func(source string, update interface{}) error

func (f MergeFunc) Merge(source string, update interface{}) error {
	// merge need to use the source message
	return f(source, update)
}

// Mux is a class for merging configuration from multiple sources.  Changes are
// pushed via channels and sent to the merge function.
type Mux struct {
	merger Merger

	// Sources and their lock.
	sourceLock sync.RWMutex

	// map from the source to the PodUpdate channel
	// use interface type here to make this class much more abstract
	sources map[string]chan interface{}
}

func NewMux(m Merger) *Mux {
	return &Mux{
		merger:  m,
		sources: map[string]chan interface{}{},
	}
}

// try to not use context to stop the job
func (m *Mux) BuildNewChan(source string) chan interface{} {
	// add a read lock
	m.sourceLock.RLock()
	ch, ok := m.sources[source]
	m.sourceLock.RUnlock()
	if ok {
		return ch
	}
	// need to build a new channel and start a new go routine to listen on it
	m.sources[source] = make(chan interface{})
	ch = m.sources[source]
	// !!!:may cause bug , use a for loop to achieve the long waiting
	go func() {
		m.listen(source, ch)
		fmt.Printf("source %s listening ending\n", source)
	}()
	return ch
}

func (m *Mux) listen(source string, ch chan interface{}) {
	// need to pay attention to close the channel, or use this method to read the message from channel may cause error
	for ele := range ch {
		m.merger.Merge(source, ele)
	}
}
