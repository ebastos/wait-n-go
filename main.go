package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Keep trying until we're timed out or got a result
func checkReady(service string, timewait time.Duration) bool {
	timeout := time.After(timewait)
	tick := time.Tick(500 * time.Millisecond)
	for {
		select {
		// Got a timeout! Stop trying and abort
		case <-timeout:
			log.Fatalf("Service %s never got reachable", service)
		// Got a tick, we should check on network connection again
		case <-tick:
			ok := checkPort(service, timewait)
			if ok {
				return true
			}

		}
	}
}

// Check if a host:port combination is listening
func checkPort(service string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", service, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func checkServices(services []string, timewait time.Duration) bool {
	var channel = make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(services))
	go func() {
		for _, service := range services {
			go func(service string) {
				defer wg.Done()
				log.Print("Waiting for ", service)
				for {
					if checkReady(service, timewait) {
						return
					}
				}
			}(service)
		}
		wg.Wait()
		close(channel)
	}()

	select {
	case <-channel: // services are ready
		log.Print("All services ready!")
		return true
	// Adding one extra second so the checkReady log.fatalf can be invoked
	case <-time.After(timewait + 1*time.Second):
		return false
	}
}

// Run the actual desired command and arguments
func startEntry(args ...string) (p *os.Process, err error) {
	log.Print("Starting entrypoint ", args)
	if args[0], err = exec.LookPath(args[0]); err == nil {
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{os.Stdin,
			os.Stdout, os.Stderr}
		p, err := os.StartProcess(args[0], args, &procAttr)
		if err == nil {
			return p, nil
		}
		fmt.Println(err)

	}
	return nil, err
}

func main() {
	var services = flag.String("services", "", "hostname to connect")
	var timeout = flag.Duration("timeout", 60*time.Second, "Maximum time to wait connection")
	flag.Parse()
	entryCommand := flag.Args()

	if *services == "" {
		fmt.Println("Missing connection information!")
		os.Exit(1)
	}
	s := strings.Split(*services, ",")

	if len(entryCommand) < 1 {
		fmt.Println("Missing entrypoint command")
		os.Exit(1)
	}

	checkServices(s, *timeout) // Will exit here is timeout.

	// Run the entry point command
	if proc, err := startEntry(entryCommand[0:]...); err == nil {
		proc.Wait()
	} else {
		fmt.Println(err)
	}

}
