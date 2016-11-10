package main

import (
	"flag"
	"io/ioutil"
	"github.com/barakb/pinger"
	"gopkg.in/tylerb/graceful.v1"
	"time"
	"log"
	"github.com/ghodss/yaml"
	"sort"
	"fmt"
)

func main() {

	cfgFilePtr := flag.String("config", "pinger.yml", "location to yaml config file")
	flag.Parse()
	config, err := readConfig(*cfgFilePtr)
	if err != nil {
		log.Fatalf("Error reading config file %v\n", err)
		panic(err)
	}
	log.Printf("config is %#v\n", config)
	pinger.InitConfig(config)
	pinger.NewmanBot.Start()
	go run(config)
	router := pinger.NewRouter()
	if err := graceful.RunWithErr(fmt.Sprintf(":%d", config.Port), 10 * time.Second, router); err != nil {
		log.Fatal(err)
	}
}

func readConfig(configFile string) (*pinger.Config, error) {
	res := &pinger.Config{}
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return res, err
	}
	return res, yaml.Unmarshal(bytes, res)
}

type ByName []pinger.Transition
func (s ByName) Len() int      { return len(s) }
func (s ByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByName) Less(i, j int) bool { return s[i].Address < s[j].Address }

func run(config *pinger.Config) {
	var total, fail int
	var transitions []pinger.Transition
	var show = false
	var oldState map[string]pinger.PingState = map[string]pinger.PingState{}
	for {
		total, fail = 0, 0
		transitions = []pinger.Transition{}
		state := pinger.RunOnce()
		total = len(state)
		//sort.Sort(pinger.PingStatusByState(state))
		for address, ps := range state {
			func(address string, ps pinger.PingStatus) {
				old, ok := oldState[address]
				if ok {
					if ps.State == pinger.Success && oldState[address] == pinger.Fail {
						show = true
						transitions = append(transitions, pinger.Transition{address, old, ps.State})
					}
					if ps.State == pinger.Fail && oldState[address] == pinger.Success {
						show = true
						transitions = append(transitions, pinger.Transition{address, old, ps.State})
					} else if ps.State == pinger.Fail {
						transitions = append(transitions, pinger.Transition{address, old, ps.State})
					}
				}else if ps.State == pinger.Fail{
					transitions = append(transitions, pinger.Transition{address, old, ps.State})
					show = true
				}
				if ps.State == pinger.Fail {
					fail += 1
				}
				oldState[address] = ps.State
			}(address, ps)
		}
		for key, _ := range oldState {
			if _, ok := state[key]; !ok {
				delete(oldState, key)

			}
		}
		if show {
			sort.Sort(ByName(transitions))
			pinger.NewmanBot.OnPingChange(transitions, total, fail)
			show = false
		}
		time.Sleep(config.WaitTime)
	}

}