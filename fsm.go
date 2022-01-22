// Copyright (c) 2022 - maozhixiang <mzx@live.cn>
// Copyright (c) 2013 - Max Persson <max@looplab.se>
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

package fsm

// FSM is the state machine model that holds the transitions and callbacks.
//
// It has to be created with NewFSM to function properly.
type FSM[STATE, EVENT comparable, FSM_IMPL, ARG any] struct {
	initial            STATE
	fsmImplConstructor func() *FSM_IMPL

	// transitions maps events and source states to destination states.
	transitions map[eKey[STATE, EVENT]]STATE

	stateCallbackFunc    map[STATE]stateCallbackFunc[STATE, EVENT, FSM_IMPL, ARG]
	eventCallbackFunc    map[EVENT]eventCallbackFunc[STATE, EVENT, FSM_IMPL, ARG]
	allStateCallbackFunc stateCallbackFunc[STATE, EVENT, FSM_IMPL, ARG]
	allEventCallbackFunc eventCallbackFunc[STATE, EVENT, FSM_IMPL, ARG]
}

// EventDesc represents an event when initializing the FSM.
//
// The event can have one or more source states that are valid for performing
// the transition. If the FSM is in one of the source states it will end up in
// the specified destination state, calling all defined callbacks as it goes.
type EventDesc[STATE, EVENT comparable] struct {
	// Name is the event name used when calling for a transition.
	Name EVENT

	// Src is a slice of source states that the FSM must be in to perform a
	// state transition.
	Src []STATE

	// Dst is the destination state that the FSM will be in if the transition
	// succeeds.
	Dst STATE
}

// Callback is a function type that callbacks should use. Event is the current
// event info as the callback happens.
type Callback[STATE, EVENT comparable, FSM_IMPL, ARG any] func(*FSM_IMPL, *Event[STATE, EVENT, ARG])

// NewFSM constructs an FSM model from events and callbacks.
//
// The events and transitions are specified as a slice of Event structs
// specified as Events. Each Event is mapped to one or more internal
// transitions from Event.Src to Event.Dst.
func NewFSM[STATE, EVENT comparable, FSM_IMPL, ARG any](initial STATE, events []EventDesc[STATE, EVENT]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	f := &FSM[STATE, EVENT, FSM_IMPL, ARG]{
		initial:           initial,
		transitions:       make(map[eKey[STATE, EVENT]]STATE),
		stateCallbackFunc: make(map[STATE]stateCallbackFunc[STATE, EVENT, FSM_IMPL, ARG]),
		eventCallbackFunc: make(map[EVENT]eventCallbackFunc[STATE, EVENT, FSM_IMPL, ARG]),
	}
	// Build transition map
	for _, e := range events {
		for _, src := range e.Src {
			f.transitions[eKey[STATE, EVENT]{e.Name, src}] = e.Dst
		}
	}
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) AddTransition(name EVENT, src []STATE, dst STATE) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	for _, src := range src {
		f.transitions[eKey[STATE, EVENT]{name, src}] = dst
	}
	return f
}

// Call method below to add callbacks at specific position of a transition.
// once transition occur, the order of callback are as follow:
//
// 1. Before - called before event named <EVENT>
//
// 2. BeforeAny - called before all events
//
// 3. OnLeave - called before leaving <STATE>
//
// 4. OnLeaveAny - called before leaving all states
//
// 5. OnEnter - called after entering <NEW_STATE>
//
// 6. OnEnterAny - called after entering all states
//
// 7. After - called after event named <EVENT>
//
// 8. AfterAny - called after all events
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) OnEnter(s STATE, cb Callback[STATE, EVENT, FSM_IMPL, ARG]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	callbackFunc := f.stateCallbackFunc[s]
	callbackFunc.enter = cb
	f.stateCallbackFunc[s] = callbackFunc
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) OnLeave(s STATE, cb Callback[STATE, EVENT, FSM_IMPL, ARG]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	callbackFunc := f.stateCallbackFunc[s]
	callbackFunc.leave = cb
	f.stateCallbackFunc[s] = callbackFunc
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) Before(e EVENT, cb Callback[STATE, EVENT, FSM_IMPL, ARG]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	callbackFunc := f.eventCallbackFunc[e]
	callbackFunc.before = cb
	f.eventCallbackFunc[e] = callbackFunc
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) After(e EVENT, cb Callback[STATE, EVENT, FSM_IMPL, ARG]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	callbackFunc := f.eventCallbackFunc[e]
	callbackFunc.after = cb
	f.eventCallbackFunc[e] = callbackFunc
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) OnEnterAny(cb Callback[STATE, EVENT, FSM_IMPL, ARG]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	f.allStateCallbackFunc.enter = cb
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) OnLeaveAny(cb Callback[STATE, EVENT, FSM_IMPL, ARG]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	f.allStateCallbackFunc.leave = cb
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) BeforeAny(cb Callback[STATE, EVENT, FSM_IMPL, ARG]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	f.allEventCallbackFunc.before = cb
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) AfterAny(cb Callback[STATE, EVENT, FSM_IMPL, ARG]) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	f.allEventCallbackFunc.after = cb
	return f
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) SetFsmImplConstructor(fsmImplConstructor func() *FSM_IMPL) *FSM[STATE, EVENT, FSM_IMPL, ARG] {
	f.fsmImplConstructor = fsmImplConstructor
	return f
}

// Can returns true if event can occur in the current state.
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) Can(current STATE, event EVENT) bool {
	_, ok := f.transitions[eKey[STATE, EVENT]{event, current}]
	return ok
}

// AvailableTransitions returns a list of transitions available in the
// current state.
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) AvailableTransitions(current STATE) []EVENT {
	var transitions []EVENT
	for key := range f.transitions {
		if key.src == current {
			transitions = append(transitions, key.event)
		}
	}
	return transitions
}

// eKey is a struct key used for storing the transition map.
type eKey[STATE, EVENT comparable] struct {
	event EVENT
	src   STATE
}

type stateCallbackFunc[STATE, EVENT comparable, FSM_IMPL, ARG any] struct {
	enter, leave Callback[STATE, EVENT, FSM_IMPL, ARG]
}
type eventCallbackFunc[STATE, EVENT comparable, FSM_IMPL, ARG any] struct {
	before, after Callback[STATE, EVENT, FSM_IMPL, ARG]
}
