package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Pin struct {
	Name string
}

type Message interface{}

type KernelFunc func([]Message, chan bool) ([]Message, error, bool)

type Kernel struct {
	Inputs   []Pin
	Outputs  []Pin
	Function KernelFunc
}

var Library = map[string]Kernel{
	"add": Kernel{
		Inputs: []Pin{
			Pin{
				Name: "addend",
			},
			Pin{
				Name: "addend",
			},
		},
		Outputs: []Pin{
			Pin{
				Name: "sum",
			},
		},
		Function: func(m []Message, quit chan bool) ([]Message, error, bool) {
			out := make([]Message, 1, 1)
			out[0] = m[0].(float64) + m[1].(float64)
			return out, nil, true
		},
	},
	"delay": Kernel{
		Inputs: []Pin{
			Pin{
				Name: "in",
			},
		},
		Outputs: []Pin{
			Pin{
				Name: "out",
			},
		},
		Function: func(m []Message, quit chan bool) ([]Message, error, bool) {
			time.Sleep(1 * time.Second)
			return m, nil, true
		},
	},
	"log": Kernel{
		Inputs: []Pin{
			Pin{
				Name: "in",
			},
		},
		Outputs: []Pin{},
		Function: func(m []Message, quit chan bool) ([]Message, error, bool) {
			fmt.Println(m[0])
			return nil, nil, true
		},
	},
	"rand": Kernel{
		Inputs: []Pin{},
		Outputs: []Pin{
			Pin{
				Name: "float",
			},
		},
		Function: func(m []Message, quit chan bool) ([]Message, error, bool) {
			out := make([]Message, 1, 1)
			out[0] = rand.Float64()
			return out, nil, true
		},
	},
}

var id int64

type KernelPin struct {
	KernelID    int64
	KernelPinID int64
}

type KernelConnection struct {
	From KernelPin
	To   KernelPin
}

type Block struct {
	Kernels           map[int64]Kernel
	KernelConnections []KernelConnection
	KernelOrder       [][]int64
}

func NewBlock() *Block {
	return &Block{
		Kernels: make(map[int64]Kernel),
	}
}

func (b *Block) AddKernel(k Kernel) int64 {
	id++
	b.Kernels[id] = k
	return id
}

func (b *Block) ConnectKernels(fromId int64, fromPinId int64, toId int64, toPinId int64) {
	nkc := KernelConnection{
		From: KernelPin{
			fromId,
			fromPinId,
		},
		To: KernelPin{
			toId,
			toPinId,
		},
	}

	b.KernelConnections = append(b.KernelConnections, nkc)
}

func (b *Block) Build() {
	lastKernels := make(map[int64]struct{})

	for k, _ := range b.Kernels {
		lastKernels[k] = struct{}{}
	}

	// search the graph for nodes that are not sending, only receiving
	// this gives us the "last" nodes in the tree
	for _, v := range b.KernelConnections {
		if _, ok := lastKernels[v.From.KernelID]; ok {
			delete(lastKernels, v.From.KernelID)
		}
	}

	var lastIDs []int64
	for k, _ := range lastKernels {
		lastIDs = append(lastIDs, k)
	}

	type orderFunc func(orderFunc, int64, int)
	type Pair struct {
		KernelID int64
		Tier     int
	}
	var tmpOrder []Pair

	// given the last node(s), search towards the top of the graph.
	// every step we take from the bottom, increment "tier"
	// each "tier" is a set of kernels that must be completed
	// for subsequent operations.
	order := func(order orderFunc, lastID int64, tier int) {
		tier++
		tmpOrder = append(tmpOrder, Pair{
			KernelID: lastID,
			Tier:     tier,
		})

		for _, v := range b.KernelConnections {
			if v.To.KernelID == lastID {
				order(order, v.From.KernelID, tier)
			}
		}
	}

	for _, v := range lastIDs {
		order(order, v, -1)
	}

	maxTier := 0
	for _, v := range tmpOrder {
		if v.Tier > maxTier {
			maxTier = v.Tier
		}
	}

	// create an 2d slice that is [tier][kernel]
	// each tier can have multiple required kernels to execute
	tierOrder := make([][]int64, maxTier+1)
	for _, v := range tmpOrder {
		if tierOrder[v.Tier] == nil {
			tierOrder[v.Tier] = []int64{}
		}
		tierOrder[v.Tier] = append(tierOrder[v.Tier], v.KernelID)
	}

	// reverse it because we started our search at the bottom of the graph
	reverse := make([][]int64, maxTier+1)
	for i, v := range tierOrder {
		reverse[maxTier-i] = v
	}

	b.KernelOrder = reverse
}

func (b *Block) Exec() {
	// horribly unoptimized exec
	// each iteration is a "tier" of kernels to be executed
	// for each loop
	//	execute each kernel on the tier
	//	use the edges to figure out where to assign the results of each kernel
	//	into out messages array
	//	this messages array is fed into the next

	// some horrible bugs just realized:
	//	1. the edges should be a map for faster lookup (duh)
	//	2. this will incorectly feed the same messages array to all kernels on a
	//	   given tier. the code that figures outbound connections needs to be added
	//	   to the beginning of the loop.
	// 	3. the messages(16) slice is a horrible apology that limits to 16 pins

	m := make([]Message, 16)
	q := make(chan bool)
	for _, t := range b.KernelOrder {
		outputs := make(map[int64][]Message)
		for _, kID := range t {
			outputs[kID], _, _ = b.Kernels[kID].Function(m, q)
		}
		m = make([]Message, 16)
		for k, v := range outputs {
			for _, c := range b.KernelConnections {
				if c.From.KernelID == k {
					m[c.To.KernelPinID] = v[c.From.KernelPinID]
				}
			}
		}
	}
}

func main() {

	nb := NewBlock()
	rId := nb.AddKernel(Library["rand"])
	r2Id := nb.AddKernel(Library["rand"])
	aId := nb.AddKernel(Library["add"])
	dId := nb.AddKernel(Library["delay"])
	lId := nb.AddKernel(Library["log"])
	nb.ConnectKernels(rId, 0, dId, 0)
	nb.ConnectKernels(dId, 0, aId, 0)
	nb.ConnectKernels(r2Id, 0, aId, 1)
	nb.ConnectKernels(aId, 0, lId, 0)
	nb.Build()
	for {
		nb.Exec()
	}
}
