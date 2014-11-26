package core

func NewPlus(name string) *Spec {
	return &Spec{
		Name:    "plus",
		Inputs:  {"addend 1", "addend 2"},
		Outputs: {"out"},
		Kernel: func(inputs map[string]interface{}) map[string]interface{} {
			output := make(map[string]interface{})
			output["out"] = values["addend 1"].(int) + values["addend 2"].(int)
			return output
		},
	}
}
