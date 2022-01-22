//go:build ignore
// +build ignore

package main

import (
	"fmt"

	fsm "github.com/maozhixiang/fsm/legacy"
)

func main() {
	var f *fsm.FSM
	f = fsm.NewFSM(
		"idle",
		fsm.Events{
			{Name: "produce", Src: []string{"idle"}, Dst: "idle"},
			{Name: "consume", Src: []string{"idle"}, Dst: "idle"},
		},
		fsm.Callbacks{
			"produce": func(e *fsm.Event) {
				f.Self.SetMetadata("message", "hii")
				fmt.Println("produced data")
			},
			"consume": func(e *fsm.Event) {
				message, ok := f.Self.Metadata("message")
				if ok {
					fmt.Println("message = " + message.(string))
				}

			},
		},
	)

	fmt.Println(f.Current())

	err := f.Event("produce")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(f.Current())

	err = f.Event("consume")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(f.Current())

}
