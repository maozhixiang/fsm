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

// Event is the info that get passed as a reference in the callbacks.
type Event[STATE, EVENT comparable, ARG any] struct {
	// Event is the event name.
	Event EVENT

	// Src is the state before the transition.
	Src STATE

	// Dst is the state after the transition.
	Dst STATE

	// Err is an optional error that can be returned from a callback.
	Err error

	// Args is an optional list of arguments passed to the callback.
	Args []ARG

	// canceled is an internal flag set if the transition is canceled.
	canceled bool
}

// Cancel can be called in before_<EVENT> or leave_<STATE> to cancel the
// current transition before it happens. It takes an optional error, which will
// overwrite e.Err if set before.
func (e *Event[STATE, EVENT, ARG]) Cancel(err ...error) {
	e.canceled = true

	if len(err) > 0 {
		e.Err = err[0]
	}
}
