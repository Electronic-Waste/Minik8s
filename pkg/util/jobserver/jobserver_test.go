package jobserver

import (
	"fmt"
	"testing"
)

func TestJobServer(t *testing.T) {
	js := NewJobServer()
	err := js.Run(true, "test.cu", "test1.slurm", "result")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("finish job server test")
}
