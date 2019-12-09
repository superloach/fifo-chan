# fifochan
Use Linux FIFO nodes as Go channels.

## Example
```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/superloach/fifochan"
)

// generic channel
var ch chan interface{}

// goroutine for reading data from fifochan
func read() {
	// forever
	for {
		// print what the channel got
		fmt.Printf("%s\n", <-ch)
	}
}

// goroutine for writing data to fifochan
func write() {
	// every .25 seconds
	for n := range time.Tick(time.Second / 4) {
		// send the time to the channel
		ch <- n
	}
}

// goroutine for reading errors from fifochan
func err() {
	// forever
	for {
		// print errors
		fmt.Printf("%s\n", <-fifochan.ErrChan())
	}
}

// main function
func main() {
	// create a new fifochan at /tmp/test
	// (assumes existence of /tmp)
	ch = fifochan.New("/tmp/test")

	// begin background streaming of data
	fifochan.Start()
	defer fifochan.Stop()

	// start the write goroutine
	go write()
	// start the read goroutine
	go read()
	// start the err goroutine
	go err()

	// catch interrupt so that defer works
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
```

## Benchmarks
```
BenchmarkChan-4         275101270              245 ns/op
BenchmarkFIFOChan-4     195536666              372 ns/op
```
