package config

import (
	"fmt"
	"testing"
	"time"
)

type testStorage struct {
}

func (t *testStorage) Merge(source string, update interface{}) error {
	// we use a int type as the channel
	fmt.Printf("get %d from %s\n", update.(int32), source)
	return nil
}

func NewTest() *testStorage {
	return &testStorage{}
}

func TestMux(t *testing.T) {
	ts := NewTest()
	mux := NewMux(ts)
	ch := mux.BuildNewChan("testSource")
	for _, i := range []int32{1, 2, 3, 4, 5} {
		ch <- i
	}
	for _, i := range []int32{6, 7, 8, 9, 10} {
		ch <- i
	}
	// need to close the channel to exit the go routine
	close(ch)
	time.Sleep(time.Second)
}
