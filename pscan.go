package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
)

var (
	hostname = flag.String("h", "localhost", "")
	cpus     = flag.Int("cpus", runtime.GOMAXPROCS(-1), "")
	poolsize = flag.Int("p", 10, "")
)
var usage = `Usage: pscan <options>

Options:
-h hostname
-cpus Number of used cpu cores. (default for current machine is %d cores)
-p no of goroutines	
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage, runtime.NumCPU()))
	}
	flag.Parse()
	runtime.GOMAXPROCS(*cpus)

	fmt.Printf("Scanning ports of %s\n", *hostname)

	// create a goroutine pool of given size. Meaning a buffered channel of that size so that
	// at a given point of time only that many are running. As we just need this to throttle goroutines
	// it can just be of type struct{} which wont consume any memory
	ch := make(chan struct{}, *poolsize)

	// to make sure we exit the program after all goroutines get done
	var wg sync.WaitGroup
	for i := 0; i < 3307; i++ {
		wg.Add(1)

		go func(port int) {
			ch <- struct{}{}
			addr := fmt.Sprintf("%s:%d", *hostname, port)
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				// handle error, just means the port is closed
				// fmt.Printf("Port %d closed\n: %v", port, err)
			} else {
				fmt.Printf("Port %d open\n", port)
				conn.Close()
			}
			// pulling out the placeholder from the channel to make space for another goroutine
			<-ch
			// marking it as complete
			wg.Done()
		}(i)
	}
	wg.Wait()
}
