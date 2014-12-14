package core

import (
	"encoding/json"
	"fmt"
	"time"
)

// return the streamtools library
func GetLibrary() map[string]Spec {
	return map[string]Spec{
		"plus": Plus(),
		"delay": Spec{
			Inputs: []Pin{
				Pin{
					"passthrough",
				},
				Pin{
					"duration",
				},
			},
			Outputs: []Pin{
				Pin{
					"passthrough",
				},
			},
			Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
				t, err := time.ParseDuration(in[1].(string))
				if err != nil {
					out[0] = err
					return nil
				}
				timer := time.NewTimer(t)
				select {
				case <-timer.C:
					out[0] = in[0]
					return nil
				case f := <-i:
					return f
				}
				return nil
			},
		},
		"set": Spec{
			Inputs: []Pin{
				Pin{
					"key",
				},
				Pin{
					"value",
				},
			},
			Outputs: []Pin{
				Pin{
					"object",
				},
			},
			Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
				out[0] = map[string]interface{}{
					in[0].(string): in[1],
				}
				return nil
			},
		},
		"log": Spec{
			Inputs: []Pin{
				Pin{
					"log",
				},
			},
			Outputs: []Pin{},
			Kernel: func(in MessageMap, out MessageMap, i chan InterruptFunc) InterruptFunc {
				o, err := json.Marshal(in[0])
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(o))
				return nil
			},
		},
	}
}
