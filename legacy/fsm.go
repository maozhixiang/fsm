// Copyright (c) 2022 - maozhixiang <mzx@live.cn>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package legacy

import (
	"strings"
	"sync"

	"github.com/maozhixiang/fsm"
)

type (
	Event                = fsm.Event[string, string, interface{}]
	Events               = []fsm.EventDesc[string, string]
	FSM                  = fsm.Instance[string, string, StateMachine, interface{}]
	Callbacks            = map[string]func(event *Event)
	InvalidEventError    = fsm.InvalidEventError[string, string]
	UnknownEventError    = fsm.UnknownEventError[string]
	NoTransitionError    = fsm.NoTransitionError
	InTransitionError    = fsm.InTransitionError[string]
	NotInTransitionError = fsm.NotInTransitionError
	CanceledError        = fsm.CanceledError
)

func NewFSM(initial string, events Events, callbacks Callbacks) *fsm.Instance[string, string, StateMachine, interface{}] {
	f := fsm.NewFSM[string, string, StateMachine, interface{}](initial, events).SetFsmImplConstructor(NewStateMachine)
	allEvents := make(map[string]bool)
	allStates := make(map[string]bool)
	for _, e := range events {
		for _, src := range e.Src {
			allStates[src] = true
			allStates[e.Dst] = true
		}
		allEvents[e.Name] = true
	}
	for name, fn := range callbacks {
		var target string
		fn := fn

		switch {
		case strings.HasPrefix(name, "before_"):
			target = strings.TrimPrefix(name, "before_")
			if target == "event" {
				f.BeforeAny(func(_ *StateMachine, e *Event) { fn(e) })
			} else if _, ok := allEvents[target]; ok {
				f.Before(target, func(_ *StateMachine, e *Event) { fn(e) })
			}
		case strings.HasPrefix(name, "leave_"):
			target = strings.TrimPrefix(name, "leave_")
			if target == "state" {
				f.OnLeaveAny(func(_ *StateMachine, e *Event) { fn(e) })
			} else if _, ok := allStates[target]; ok {
				f.OnLeave(target, func(_ *StateMachine, e *Event) { fn(e) })
			}
		case strings.HasPrefix(name, "enter_"):
			target = strings.TrimPrefix(name, "enter_")
			if target == "state" {
				f.OnEnterAny(func(_ *StateMachine, e *Event) { fn(e) })
			} else if _, ok := allStates[target]; ok {
				f.OnEnter(target, func(_ *StateMachine, e *Event) { fn(e) })
			}
		case strings.HasPrefix(name, "after_"):
			target = strings.TrimPrefix(name, "after_")
			if target == "event" {
				f.AfterAny(func(_ *StateMachine, e *Event) { fn(e) })
			} else if _, ok := allEvents[target]; ok {
				f.After(target, func(_ *StateMachine, e *Event) { fn(e) })
			}
		default:
			target = name
			if _, ok := allStates[target]; ok {
				f.OnEnter(target, func(_ *StateMachine, e *Event) { fn(e) })
			} else if _, ok := allEvents[target]; ok {
				f.After(target, func(_ *StateMachine, e *Event) { fn(e) })
			}
		}
	}
	return f.NewInstance()
}

type StateMachine struct {
	// metadata can be used to store and load data that maybe used across events
	// use methods SetMetadata() and Metadata() to store and load data
	metadata map[string]interface{}

	metadataMu sync.RWMutex
}

func NewStateMachine() *StateMachine {
	return &StateMachine{metadata: map[string]interface{}{}}
}

// Metadata returns the value stored in metadata
func (f *StateMachine) Metadata(key string) (interface{}, bool) {
	f.metadataMu.RLock()
	defer f.metadataMu.RUnlock()
	dataElement, ok := f.metadata[key]
	return dataElement, ok
}

// SetMetadata stores the dataValue in metadata indexing it with key
func (f *StateMachine) SetMetadata(key string, dataValue interface{}) {
	f.metadataMu.Lock()
	defer f.metadataMu.Unlock()
	f.metadata[key] = dataValue
}
