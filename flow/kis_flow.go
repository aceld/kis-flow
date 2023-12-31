package flow

import "kis-flow/config"

// KisFlow 用于贯穿整条流式计算的上下文环境
type KisFlow struct {
	Id   string
	Name string
	// TODO
}

// TODO for test
func NewKisFlow(config *config.KisFlowConfig) *KisFlow {
	flow := new(KisFlow)
	flow.Id = config.FlowId
	flow.Name = config.FlowName

	return flow
}
