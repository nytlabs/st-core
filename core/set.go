package core

func SetSpec() Spec {
	return Spec{
		Name:    "set",
		Inputs:  []string{"in", "key"},
		Outputs: []string{"out"},
		Kernel:  SetKernel,
	}
}

func SetKernel(quit chan bool, inputs map[string]Message) (map[string]Message, bool) {

	out := make(map[string]Message)
	k := inputs["key"].(string)
	out["out"] = map[string]interface{}{
		k: inputs["in"],
	}

	return out, true
}
