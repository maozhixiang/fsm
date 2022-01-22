# FSM for Go

FSM is a finite state machine for Go.

Forked from github.com/looplab/fsm.

Different from the previous implementation:  
- for better memory allocation, instances created from one FSM model will share the definition of transition and callback
- using type parameters for state, event, fsm_struct and event_arg. (requiring go 1.18).
- a 'legacy' package is provided for previous implementation and test
- no longer provide metadata or asynchronous state transition, witch can be implemented in custom FSM struct

# Basic Example

From examples/simple.go:

```go
package main

import (
    "fmt"

    fsm "github.com/maozhixiang/fsm/legacy"
)

func main() {
    fsm := fsm.NewFSM(
        "closed",
        fsm.Events{
            {Name: "open", Src: []string{"closed"}, Dst: "open"},
            {Name: "close", Src: []string{"open"}, Dst: "closed"},
        },
        fsm.Callbacks{},
    )

    fmt.Println(fsm.Current())

    err := fsm.Event("open")
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(fsm.Current())

    err = fsm.Event("close")
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(fsm.Current())
}
```

# Usage as a struct field

From examples/struct.go:

```go
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
```

# License

FSM is licensed under Apache License 2.0

http://www.apache.org/licenses/LICENSE-2.0
