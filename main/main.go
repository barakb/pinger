package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	pinger "github.com/barakb/pinger"
)

type Config struct {
	Machines []string `json:"machines"`
}

func main() {
	cfgFilePtr := flag.String("config", "pinger.json", "location to cfg file")
	flag.Parse()
	config, err := readConfig(*cfgFilePtr)
	if err != nil {
		fmt.Printf("Error reading config file %v\n", err)
		panic(err)
	}
	fmt.Printf("config is %#v\n", config)
	for {
		pingsResult := pinger.ParallelPing(config.Machines)
		for _, pingResult := range pingsResult {
			// if pingResult.Err != nil {
			// 	fmt.Printf("error: %v while pinging to %q\n", pingResult.Err, pingResult.Machine)
			// }
			fmt.Printf("ping to %q with status: %d, elapsed time: %v, error is: %v\n", pingResult.Machine, pingResult.Status, pingResult.Elapsed, pingResult.Err)
		}
		time.Sleep(1 * time.Second)
	}
}

func readConfig(configFile string) (config *Config, err error) {
	res := &Config{}
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}
	return res, json.Unmarshal(bytes, res)
}
