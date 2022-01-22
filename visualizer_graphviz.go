package fsm

import (
	"bytes"
	"fmt"
)

// Visualize outputs a visualization of a FSM in Graphviz format.
func (fsm *FSM[STATE, EVENT, FSM_IMPL, ARG]) Visualize(current STATE) string {
	var buf bytes.Buffer

	// we sort the key alphabetically to have a reproducible graph output
	sortedEKeys := fsm.getSortedTransitionKeys()
	sortedStateKeys, _ := fsm.getSortedStates()

	//writeHeaderLine(&buf)
	buf.WriteString(fmt.Sprintf(`digraph fsm {`))
	buf.WriteString("\n")

	//writeTransitions(&buf, fmt.Sprint(current), sortedEKeys, fsm.transitions)
	for _, k := range sortedEKeys {
		if k.src == current {
			v := fsm.transitions[k]
			buf.WriteString(fmt.Sprintf(`    "%v" -> "%v" [ label = "%v" ];`, k.src, v, k.event))
			buf.WriteString("\n")
		}
	}
	for _, k := range sortedEKeys {
		if k.src != current {
			v := fsm.transitions[k]
			buf.WriteString(fmt.Sprintf(`    "%v" -> "%v" [ label = "%v" ];`, k.src, v, k.event))
			buf.WriteString("\n")
		}
	}

	// make sure the current state is at top
	buf.WriteString("\n")

	//writeStates(&buf, sortedStateKeys)
	for _, k := range sortedStateKeys {
		buf.WriteString(fmt.Sprintf(`    "%v";`, k))
		buf.WriteString("\n")
	}

	//writeFooter(&buf)
	buf.WriteString(fmt.Sprintln("}"))

	return buf.String()
}
