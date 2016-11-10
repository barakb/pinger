package pinger

import (
	"os/exec"
	"sync"
	"syscall"
	"time"
	"errors"
	"bytes"
	"log"
)

type PingResult struct {
	Status    int `json:"status"`
	Elapsed   time.Duration `json:"elapsed"`
	Err       error `json:"error"`
	Address   string `json:"address"`
	Timestamp time.Time `json:"timestamp"`
}

type State struct {
	history  map[string][]*PingResult
	statuses map[string]PingStatus
}

func (state *State) updatePingsResults(pingResult *PingResult) {
	history, ok := state.history[pingResult.Address]
	if !ok {
		history = []*PingResult{}
		state.history[pingResult.Address] = history
	}
	history = append(history, pingResult)
	if 5 < len(history) {
		history = history[1:]
	}
	state.history[pingResult.Address] = history
	success, fail, avg := computeState(history)
	//log.Printf("* %s: sucessu: %d, fail: %d, avg: %s\n", pingResult.Address, success, fail, avg)

	if fail < success || fail == 0{
		state.statuses[pingResult.Address] = PingStatus{Success, pingResult.Address, pingResult.Timestamp, avg}
	}else{
		state.statuses[pingResult.Address] = PingStatus{Fail, pingResult.Address, pingResult.Timestamp, avg}
	}
}

func computeState(results []*PingResult) (success, fail int, avg time.Duration) {
	var nanos int64
	for _, result := range results {
		if result.Err != nil || result.Status != 0 {
			fail += 1
		} else {
			success += 1
			nanos += result.Elapsed.Nanoseconds()
		}
	}
	if 0 < success {
		avg = time.Duration(nanos / int64(success))
	}
	return;
}

func (state *State) updateStatus() {
	for _, address := range config.Addresses {
		history, ok := state.history[address]
		if !ok {
			history = []*PingResult{}
			state.history[address] = history
		}

	}
}

type PingStatus struct {
	State              PingState`json:"state"`
	Address            string `json:"address"`
	Timestamp          time.Time `json:"timestamp"`
	AverageSuccessTime time.Duration `json:"average_sucess_time"`
}


type PingState int

func (ps PingState) String() string {
	switch ps{
	case Success : return "success"
	case Fail : return "fail"
	default:
		return "Unknown"
	}
}

const (
	Success PingState = 1
	Fail PingState = 2
)

type Transition struct {
	Address  string
	From, To PingState
}

var state *State

func init() {
	state = &State{map[string][]*PingResult{}, map[string]PingStatus{}}
}

func ping(address string, timeout time.Duration) *PingResult {
	var waitStatus syscall.WaitStatus
	var status int
	var elapsed time.Duration

	start := time.Now()
	cmd := exec.Command("ping", "-c", "1", address)

	//b, err := CombinedOutput(cmd, timeout)
	_, err := CombinedOutput(cmd, timeout)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			status = waitStatus.ExitStatus()
			//log.Printf("************** Error %s while pinging to %q, output is: %s\n", err.Error(), address, string(b[:]))
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		status = waitStatus.ExitStatus()
	}
	elapsed = time.Since(start)
	log.Printf("ping to %q took %s\n", address, elapsed)
	return &PingResult{status, elapsed, err, address, time.Now()}
}

func ParallelPing(addresses []string, timeout time.Duration) (pingsResult []*PingResult) {
	var wg sync.WaitGroup
	wg.Add(len(addresses))
	c := make(chan *PingResult, len(addresses))
	for _, address := range addresses {
		go func(address string) {
			defer wg.Done()
			c <- ping(address, timeout)
		}(address)
	}

	go func() {
		wg.Wait()
		close(c)
	}()
	now := time.Now()
	for pingResult := range c {
		pingResult.Timestamp = now
		pingsResult = append(pingsResult, pingResult)
	}
	return pingsResult
}

func CombinedOutput(cmd *exec.Cmd, timeout time.Duration) ([]byte, error) {
	if cmd.Stdout != nil {
		return nil, errors.New("exec: Stdout already set")
	}
	if cmd.Stderr != nil {
		return nil, errors.New("exec: Stderr already set")
	}
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	if err := cmd.Start(); err != nil {
		return b.Bytes(), err
	}
	timer := time.AfterFunc(timeout, func() {
		cmd.Process.Kill()
	})
	defer timer.Stop()
	return b.Bytes(), cmd.Wait()
}

func RunOnce() map[string]PingStatus {
	pingsResult := ParallelPing(config.Addresses, config.Timeout)
	for _, pingResult := range pingsResult {
		state.updatePingsResults(pingResult)
	}
	state.updateStatus()
	return state.statuses
}


