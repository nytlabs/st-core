package core

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var Library = map[string]Spec{
	"plus": Spec{
		Name:    "plus",
		Inputs:  []string{"addend 1", "addend 2"},
		Outputs: []string{"out"},
		Kernel: func(quit chan bool, inputs map[string]Message) (map[string]Message, bool) {
			output := make(map[string]Message)
			output["out"] = inputs["addend 1"].(int) + inputs["addend 2"].(int)
			return output, true
		},
	},
	"log": Spec{
		Name:    "log",
		Inputs:  []string{"in"},
		Outputs: []string{},
		Kernel: func(quit chan bool, inputs map[string]Message) (map[string]Message, bool) {
			fmt.Println(inputs["in"])
			return nil, true
		},
	},
	"delay": Spec{
		Name:    "log",
		Inputs:  []string{"in", "duration"},
		Outputs: []string{"out"},
		Kernel: func(quit chan bool, inputs map[string]Message) (map[string]Message, bool) {
			output := make(map[string]Message)
			output["out"] = inputs["in"]
			durationString, ok := inputs["duration"].(string)
			if !ok {
				log.Fatal("could not assert duration to string")
			}
			d, err := time.ParseDuration(durationString)
			if err != nil {
				log.Fatal("could not parse duration string")
			}
			t := time.NewTimer(d)
			select {
			case <-quit:
				return nil, false
			case <-t.C:
				break
			}
			return output, true
		},
	},
	"pusher": Spec{
		Name:    "pusher",
		Inputs:  []string{"value"},
		Outputs: []string{"out"},
		Kernel: func(quit chan bool, inputs map[string]Message) (map[string]Message, bool) {
			output := make(map[string]Message)
			output["out"] = inputs["value"]
			return output, true
		},
	},
	"latch": Spec{
		Name:    "latch",
		Inputs:  []string{"in", "ctrl"},
		Outputs: []string{"out"},
		Kernel: func(quit chan bool, inputs map[string]Message) (map[string]Message, bool) {
			ctrlSignal := inputs["ctrl"]
			output := make(map[string]Message)
			switch ctrlSignal := ctrlSignal.(type) {
			case bool:
				if ctrlSignal {
					output["out"] = inputs["in"]
					return output, true
				}
			case error:
				log.Fatal(ctrlSignal)
			default:
				log.Fatal(errors.New("unrecognised control signal in latch"))
			}
			return nil, true
		},
	},
	"set": Set(),
}
