package machines

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type PingResult struct {
	Status  int
	Elapsed time.Duration
	Err     error
	Machine string
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func ping(address string) (pingResult *PingResult) {
	var waitStatus syscall.WaitStatus
	var status int
	var elapsed time.Duration
	var err error

	start := time.Now()
	// Create an *exec.Cmd
	cmd := exec.Command("ping", "-c", "1", address)

	// Combine stdout and stderr
	// printCommand(cmd)
	_, err = cmd.CombinedOutput()
	if err != nil {
		// printError(err)
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			status = waitStatus.ExitStatus()
			//printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		//printOutput(output)
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		status = waitStatus.ExitStatus()
		//printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
	elapsed = time.Since(start)
	return &PingResult{status, elapsed, err, address}
}

func ParallelPing(addresses []string) (pingsResult []*PingResult) {
	var wg sync.WaitGroup
	wg.Add(len(addresses))
	c := make(chan *PingResult, len(addresses))
	for _, address := range addresses {
		go func(address string) {
			defer wg.Done()
			c <- ping(address)
		}(address)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for pingResult := range c {
		pingsResult = append(pingsResult, pingResult)
	}
	return pingsResult
}
