package fifochan

import (
	"testing"
)

const data string = "owo"

func BenchmarkChan(b *testing.B) {
	bar := Make()
	go func() {
		for i := 0; i < b.N; i++ {
			bar <- data
		}
	}()
	for i := 0; i < b.N; i++ {
		<-bar
	}
}

func BenchmarkFIFOChan(b *testing.B) {
	foo := New("/tmp/test")
	Start()
	defer Stop()
	go func() {
		for i := 0; i < b.N; i++ {
			foo <- data
		}
	}()
	for i := 0; i < b.N; i++ {
		<-foo
	}
}
