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
			{Name: "scan", Src: []string{"idle"}, Dst: "scanning"},
			{Name: "working", Src: []string{"scanning"}, Dst: "scanning"},
			{Name: "situation", Src: []string{"scanning"}, Dst: "scanning"},
			{Name: "situation", Src: []string{"idle"}, Dst: "idle"},
			{Name: "finish", Src: []string{"scanning"}, Dst: "idle"},
		},
		fsm.Callbacks{
			"scan": func(e *fsm.Event) {
				fmt.Println("after_scan: " + f.Current())
			},
			"working": func(e *fsm.Event) {
				fmt.Println("working: " + f.Current())
			},
			"situation": func(e *fsm.Event) {
				fmt.Println("situation: " + f.Current())
			},
			"finish": func(e *fsm.Event) {
				fmt.Println("finish: " + f.Current())
			},
		},
	)

	fmt.Println(f.Current())

	err := f.Event("scan")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("1:" + f.Current())

	err = f.Event("working")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("2:" + f.Current())

	err = f.Event("situation")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("3:" + f.Current())

	err = f.Event("finish")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("4:" + f.Current())

}
