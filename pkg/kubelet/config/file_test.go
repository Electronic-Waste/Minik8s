package config

import (
	"testing"
)

func TestNewFile(t *testing.T) {
	ch := make(chan interface{})
	NewSourceFile(ch)
	<-ch
}
