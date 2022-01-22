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

package fsm

import (
	"sync"
)

type Instance[STATE, EVENT comparable, FSM_IMPL, ARG any] struct {
	*FSM[STATE, EVENT, FSM_IMPL, ARG]
	Self *FSM_IMPL
	// current is the state that the FSM is currently in.
	current STATE

	// stateMu guards access to the current state.
	stateMu sync.RWMutex
	// eventMu guards access to Event() and Transition().
	eventMu sync.Mutex
}

func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) NewInstance() *Instance[STATE, EVENT, FSM_IMPL, ARG] {
	var impl *FSM_IMPL
	if f.fsmImplConstructor != nil {
		impl = f.fsmImplConstructor()
	} else {
		impl = new(FSM_IMPL)
	}
	return &Instance[STATE, EVENT, FSM_IMPL, ARG]{
		FSM:     f,
		Self:    impl,
		current: f.initial,
	}
}
func (f *FSM[STATE, EVENT, FSM_IMPL, ARG]) NewInstanceWithImpl(impl *FSM_IMPL) *Instance[STATE, EVENT, FSM_IMPL, ARG] {
	return &Instance[STATE, EVENT, FSM_IMPL, ARG]{
		FSM:     f,
		Self:    impl,
		current: f.initial,
	}
}

// Current returns the current state of the FSM.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) Current() STATE {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.current
}

// Is returns true if state is the current state.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) Is(state STATE) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return state == f.current
}

// SetState allows the user to move to the given state from current state.
// The call does not trigger any callbacks, if defined.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) SetState(state STATE) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()
	f.current = state
	return
}

// Can returns true if event can occur in the current state.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) Can(event EVENT) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.FSM.Can(f.current, event)
}

// AvailableTransitions returns a list of transitions available in the
// current state.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) AvailableTransitions() []EVENT {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()
	return f.FSM.AvailableTransitions(f.current)
}

// Cannot returns true if event can not occur in the current state.
// It is a convenience method to help code read nicely.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) Cannot(event EVENT) bool {
	return !f.Can(event)
}

// Event initiates a state transition with the named event.
//
// The call takes a variable number of arguments that will be passed to the
// callback, if defined.
//
// It will return nil if the state change is ok or one of these errors:
//
// - event X inappropriate in current state Y
//
// - event X does not exist
//
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) Event(event EVENT, args ...ARG) error {
	f.eventMu.Lock()
	defer f.eventMu.Unlock()

	f.stateMu.RLock()
	defer f.stateMu.RUnlock()

	dst, ok := f.transitions[eKey[STATE, EVENT]{event, f.current}]
	if !ok {
		for ekey := range f.transitions {
			if ekey.event == event {
				return InvalidEventError[STATE, EVENT]{event, f.current}
			}
		}
		return UnknownEventError[EVENT]{event}
	}

	e := &Event[STATE, EVENT, ARG]{event, f.current, dst, nil, args, false}

	err := f.beforeEventCallbacks(e)
	if err != nil {
		return err
	}

	if f.current == dst {
		f.afterEventCallbacks(e)
		return NoTransitionError{e.Err}
	}

	if err = f.leaveStateCallbacks(e); err != nil {
		return err
	}

	f.current = dst
	f.enterStateCallbacks(e)
	f.afterEventCallbacks(e)

	return e.Err
}

// beforeEventCallbacks calls the before_ callbacks, first the named then the
// general version.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) beforeEventCallbacks(e *Event[STATE, EVENT, ARG]) error {
	if fn := f.eventCallbackFunc[e.Event].before; fn != nil {
		fn(f.Self, e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	if fn := f.allEventCallbackFunc.before; fn != nil {
		fn(f.Self, e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	return nil
}

// leaveStateCallbacks calls the leave_ callbacks, first the named then the
// general version.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) leaveStateCallbacks(e *Event[STATE, EVENT, ARG]) error {
	if fn := f.stateCallbackFunc[f.current].leave; fn != nil {
		fn(f.Self, e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	if fn := f.allStateCallbackFunc.leave; fn != nil {
		fn(f.Self, e)
		if e.canceled {
			return CanceledError{e.Err}
		}
	}
	return nil
}

// enterStateCallbacks calls the enter_ callbacks, first the named then the
// general version.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) enterStateCallbacks(e *Event[STATE, EVENT, ARG]) {
	if fn := f.stateCallbackFunc[f.current].enter; fn != nil {
		fn(f.Self, e)
	}
	if fn := f.allStateCallbackFunc.enter; fn != nil {
		fn(f.Self, e)
	}
}

// afterEventCallbacks calls the after_ callbacks, first the named then the
// general version.
func (f *Instance[STATE, EVENT, FSM_IMPL, ARG]) afterEventCallbacks(e *Event[STATE, EVENT, ARG]) {
	if fn := f.eventCallbackFunc[e.Event].after; fn != nil {
		fn(f.Self, e)
	}
	if fn := f.allEventCallbackFunc.after; fn != nil {
		fn(f.Self, e)
	}
}
