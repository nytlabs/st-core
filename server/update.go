package server

import (
	"time"

	"github.com/nytlabs/st-core/core"
)

type Update struct {
	Action string      `json:"action"`
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
}

type wsId struct {
	Id int `json:"id"`
}

type wsLabel struct {
	wsId
	Label string `json:"label"`
}

type wsPosition struct {
	wsId
	Position Position `json:"position"`
}

// type BLOCK
type wsBlock struct {
	Block interface{} `json:"block"`
}

// type GROUP
type wsGroup struct {
	Group interface{} `json:"group"`
}

// type SOURCE
type wsSource struct {
	Source interface{} `json:"source"`
}

// type LINK
type wsLink struct {
	Link interface{} `json:"link"`
}

// type CONNECTION
type wsConnection struct {
	Connection interface{} `json:"connection"`
}

// type CHILD
type wsGroupChild struct {
	Group wsId `json:"group"`
	Child wsId `json:"child"`
}

// type ROUTE
type wsRouteModify struct {
	ConnectionNode
	Value *core.InputValue `json:"value"`
}

// type PARAM
type wsSourceModify struct {
	wsId
	Param string `json:"param"`
	Value string `json:"value"`
}

//type ALERT
type wsAlert struct {
	wsId
	Alert     core.BlockAlert `json:"alert"`
	LastCrank time.Time       `json:"lastCrank"`
}
