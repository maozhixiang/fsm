package main

import (
	"fmt"
	"strings"

	"github.com/maozhixiang/fsm"
)

type State int8
type Event int8

const (
	closed State = iota
	openState
)
const (
	open Event = iota
	close
)

type Door struct {
	To string
}

func NewDoor(to string) *Door { return &Door{To: to} }

func (d *Door) enterState(e *fsm.Event[State, Event, string]) {
	fmt.Printf("The door to %v is %v args:%s \n", d.To, e.Dst, strings.Join(e.Args, ","))
}

var doorFSM = fsm.NewFSM[State, Event, Door, string](closed, nil).
	AddTransition(open, []State{closed}, openState).
	AddTransition(close, []State{openState}, closed).
	OnEnterAny((*Door).enterState)

func main() {
	door := doorFSM.NewInstanceWithImpl(NewDoor("homeland"))

	err := door.Event(open, "args1", "args2")
	if err != nil {
		fmt.Println(err)
	}

	err = door.Event(close, "args3", "args4", "args5", "args6")
	if err != nil {
		fmt.Println(err)
	}
}
