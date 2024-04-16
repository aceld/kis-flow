package flow

import (
	"context"
	"errors"
	"fmt"

	"github.com/aceld/kis-flow/kis"
)

// dealAction handles Action to determine the next direction of the Flow.
func (flow *KisFlow) dealAction(ctx context.Context, fn kis.Function) (kis.Function, error) {

	// DataReuse Action
	if flow.action.DataReuse {
		if err := flow.commitReuseData(ctx); err != nil {
			return nil, err
		}
	} else {
		if err := flow.commitCurData(ctx); err != nil {
			return nil, err
		}
	}

	// ForceEntryNext Action
	if flow.action.ForceEntryNext {
		if err := flow.commitVoidData(ctx); err != nil {
			return nil, err
		}
		flow.abort = false
	}

	// JumpFunc Action
	if flow.action.JumpFunc != "" {
		if _, ok := flow.Funcs[flow.action.JumpFunc]; !ok {
			// The current JumpFunc is not in the flow
			return nil, errors.New(fmt.Sprintf("Flow Jump -> %s is not in Flow", flow.action.JumpFunc))
		}

		jumpFunction := flow.Funcs[flow.action.JumpFunc]
		// Update the upper layer Function
		flow.PrevFunctionId = jumpFunction.GetPrevId()
		fn = jumpFunction

		// If set to jump, force the jump
		flow.abort = false

	} else {

		// Update the upper layer FunctionId cursor
		flow.PrevFunctionId = flow.ThisFunctionId
		fn = fn.Next()
	}

	// Abort Action force termination
	if flow.action.Abort {
		flow.abort = true
	}

	// Clear Action
	flow.action = kis.Action{}

	return fn, nil
}
