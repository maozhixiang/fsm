package fsm

import (
	"bytes"
	"fmt"
)

const highlightingColor = "#00AA00"

// MermaidDiagramType the type of the mermaid diagram type
type MermaidDiagramType string

const (
	// FlowChart the diagram type for output in flowchart style (https://mermaid-js.github.io/mermaid/#/flowchart) (including current state)
	FlowChart MermaidDiagramType = "flowChart"
	// StateDiagram the diagram type for output in stateDiagram style (https://mermaid-js.github.io/mermaid/#/stateDiagram)
	StateDiagram MermaidDiagramType = "stateDiagram"
)

// VisualizeForMermaidWithGraphType outputs a visualization of a FSM in Mermaid format as specified by the graphType.
func (fsm *FSM[STATE, EVENT, FSM_IMPL, ARG]) VisualizeForMermaidWithGraphType(graphType MermaidDiagramType, current STATE) (string, error) {
	switch graphType {
	case FlowChart:
		return fsm.visualizeForMermaidAsFlowChart(current), nil
	case StateDiagram:
		return fsm.visualizeForMermaidAsStateDiagram(current), nil
	default:
		return "", fmt.Errorf("unknown MermaidDiagramType: %s", graphType)
	}
}

func (fsm *FSM[STATE, EVENT, FSM_IMPL, ARG]) visualizeForMermaidAsStateDiagram(current STATE) string {
	var buf bytes.Buffer

	sortedTransitionKeys := fsm.getSortedTransitionKeys()

	buf.WriteString("stateDiagram-v2\n")
	buf.WriteString(fmt.Sprintln(`    [*] -->`, current))

	for _, k := range sortedTransitionKeys {
		v := fsm.transitions[k]
		buf.WriteString(fmt.Sprintf(`    %s --> %s: %s`, k.src, v, k.event))
		buf.WriteString("\n")
	}

	return buf.String()
}

// visualizeForMermaidAsFlowChart outputs a visualization of a FSM in Mermaid format (including highlighting of current state).
func (fsm *FSM[STATE, EVENT, FSM_IMPL, ARG]) visualizeForMermaidAsFlowChart(current STATE) string {
	var buf bytes.Buffer

	sortedTransitionKeys := fsm.getSortedTransitionKeys()
	sortedStates, statesToIDMap := fsm.getSortedStates()

	//writeFlowChartGraphType(&buf)
	buf.WriteString("graph LR\n")

	//writeFlowChartStates(&buf, sortedStates, statesToIDMap)
	for _, state := range sortedStates {
		buf.WriteString(fmt.Sprintf(`    %s[%v]`, statesToIDMap[state], state))
		buf.WriteString("\n")
	}
	buf.WriteString("\n")

	//writeFlowChartTransitions(&buf, fsm.transitions, sortedTransitionKeys, statesToIDMap)
	for _, transition := range sortedTransitionKeys {
		target := fsm.transitions[transition]
		buf.WriteString(fmt.Sprintf(`    %s --> |%v| %v`, statesToIDMap[transition.src], transition.event, statesToIDMap[target]))
		buf.WriteString("\n")
	}
	buf.WriteString("\n")

	//writeFlowChartHighlightCurrent(&buf, fsm.current, statesToIDMap)
	buf.WriteString(fmt.Sprintf(`    style %s fill:%v`, statesToIDMap[current], highlightingColor))
	buf.WriteString("\n")

	return buf.String()
}
