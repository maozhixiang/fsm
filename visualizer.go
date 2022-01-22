package fsm

import (
	"fmt"
	"sort"
)

// VisualizeType the type of the visualization
type VisualizeType string

const (
	// GRAPHVIZ the type for graphviz output (http://www.webgraphviz.com/)
	GRAPHVIZ VisualizeType = "graphviz"
	// MERMAID the type for mermaid output (https://mermaid-js.github.io/mermaid/#/stateDiagram) in the stateDiagram form
	MERMAID VisualizeType = "mermaid"
	// MermaidStateDiagram the type for mermaid output (https://mermaid-js.github.io/mermaid/#/stateDiagram) in the stateDiagram form
	MermaidStateDiagram VisualizeType = "mermaid-state-diagram"
	// MermaidFlowChart the type for mermaid output (https://mermaid-js.github.io/mermaid/#/flowchart) in the flow chart form
	MermaidFlowChart VisualizeType = "mermaid-flow-chart"
)

// VisualizeWithType outputs a visualization of a FSM in the desired format.
// If the type is not given it defaults to GRAPHVIZ
func (fsm *FSM[STATE, EVENT, FSM_IMPL, ARG]) VisualizeWithType(visualizeType VisualizeType, current STATE) (string, error) {
	switch visualizeType {
	case GRAPHVIZ:
		return fsm.Visualize(current), nil
	case MERMAID, MermaidStateDiagram:
		return fsm.VisualizeForMermaidWithGraphType(StateDiagram, current)
	case MermaidFlowChart:
		return fsm.VisualizeForMermaidWithGraphType(FlowChart, current)
	default:
		return "", fmt.Errorf("unknown VisualizeType: %s", visualizeType)
	}
}

func (fsm *Instance[STATE, EVENT, FSM_IMPL, ARG]) VisualizeWithType(visualizeType VisualizeType) (string, error) {
	return fsm.FSM.VisualizeWithType(visualizeType, fsm.Current())
}

func (fsm *FSM[STATE, EVENT, FSM_IMPL, ARG]) getSortedTransitionKeys() []eKey[STATE, EVENT] {
	// we sort the key alphabetically to have a reproducible graph output
	sortedTransitionKeys := make([]eKey[STATE, EVENT], 0)

	for transition := range fsm.transitions {
		sortedTransitionKeys = append(sortedTransitionKeys, transition)
	}

	sort.Slice(sortedTransitionKeys, func(i, j int) bool {
		if sortedTransitionKeys[i].src == sortedTransitionKeys[j].src {
			return fmt.Sprint(sortedTransitionKeys[i].event) < fmt.Sprint(sortedTransitionKeys[j].event)
		}
		return fmt.Sprint(sortedTransitionKeys[i].src) < fmt.Sprint(sortedTransitionKeys[j].src)
	})

	return sortedTransitionKeys
}

func (fsm *FSM[STATE, EVENT, FSM_IMPL, ARG]) getSortedStates() ([]STATE, map[STATE]string) {
	statesToIDMap := make(map[STATE]string)
	for transition, target := range fsm.transitions {
		if _, ok := statesToIDMap[transition.src]; !ok {
			statesToIDMap[transition.src] = ""
		}
		if _, ok := statesToIDMap[target]; !ok {
			statesToIDMap[target] = ""
		}
	}

	sortedStates := make([]STATE, 0, len(statesToIDMap))
	for state := range statesToIDMap {
		sortedStates = append(sortedStates, state)
	}
	sort.Slice(sortedStates, func(i, j int) bool {
		return fmt.Sprint(sortedStates[i]) < fmt.Sprint(sortedStates[j])
	})

	for i, state := range sortedStates {
		statesToIDMap[state] = fmt.Sprintf("id%d", i)
	}
	return sortedStates, statesToIDMap
}
