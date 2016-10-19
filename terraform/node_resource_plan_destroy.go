package terraform

import (
	"fmt"
)

// NodePlanDestroyableResource represents a resource that is "applyable":
// it is ready to be applied and is represented by a diff.
type NodePlanDestroyableResource struct {
	*NodeAbstractResource
}

// GraphNodeEvalable
func (n *NodePlanDestroyableResource) EvalTree() EvalNode {
	addr := n.NodeAbstractResource.Addr

	// stateId is the ID to put into the state
	stateId := addr.stateId()
	if addr.Index > -1 {
		stateId = fmt.Sprintf("%s.%d", stateId, addr.Index)
	}

	// Build the instance info. More of this will be populated during eval
	info := &InstanceInfo{
		Id:   stateId,
		Type: addr.Type,
	}

	// Declare a bunch of variables that are used for state during
	// evaluation. Most of this are written to by-address below.
	var diff *InstanceDiff
	var state *InstanceState

	return &EvalSequence{
		Nodes: []EvalNode{
			&EvalReadState{
				Name:   stateId,
				Output: &state,
			},
			&EvalDiffDestroy{
				Info:   info,
				State:  &state,
				Output: &diff,
			},
			&EvalCheckPreventDestroy{
				Resource: n.Config,
				Diff:     &diff,
			},
			&EvalWriteDiff{
				Name: stateId,
				Diff: &diff,
			},
		},
	}
}
