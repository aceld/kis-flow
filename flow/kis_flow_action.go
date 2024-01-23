package flow

import (
	"context"
	"errors"
	"fmt"
	"kis-flow/kis"
)

// dealAction  处理Action，决定接下来Flow的流程走向
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
			//当前JumpFunc不在flow中
			return nil, errors.New(fmt.Sprintf("Flow Jump -> %s is not in Flow", flow.action.JumpFunc))
		}

		jumpFunction := flow.Funcs[flow.action.JumpFunc]
		// 更新上层Function
		flow.PrevFunctionId = jumpFunction.GetPrevId()
		fn = jumpFunction

		// 如果设置跳跃，强制跳跃
		flow.abort = false

	} else {

		// 更新上一层 FuncitonId 游标
		flow.PrevFunctionId = flow.ThisFunctionId
		fn = fn.Next()
	}

	// Abort Action 强制终止
	if flow.action.Abort {
		flow.abort = true
	}

	// 清空Action
	flow.action = kis.Action{}

	return fn, nil
}
