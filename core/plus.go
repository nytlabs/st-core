package core

import (
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
		Inputs:  []string{"in"},
		Outputs: []string{"out"},
		Kernel: func(quit chan bool, inputs map[string]Message) (map[string]Message, bool) {
			log.Println("DELAY", inputs)
			output := make(map[string]Message)

			fmt.Println("delay", inputs["in"])

			output["out"] = inputs["in"]
			t := time.NewTimer(1 * time.Second)
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
		Inputs:  []string{},
		Outputs: []string{"out"},
		Kernel: func(quit chan bool, inputs map[string]Message) (map[string]Message, bool) {
			output := make(map[string]Message)
			output["out"] = "ello!"
			return output, true
		},
	},
}
