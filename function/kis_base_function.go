package function

import (
	"context"
	"errors"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
	"kis-flow/id"
)

type BaseFunction struct {
	Config *config.KisFuncConfig

	Flow *flow.KisFlow //上下文环境KisFlow
	cid  string        //当前Function所依赖的KisConnectorID(如果存在)

	N KisFunction //下一个流计算Function
	P KisFunction //上一个流计算Function

	//KisId , KisFunction的实例ID，用于KisFlow内部区分不同的实例对象
	//KisId 和 Function Config中的 Fid的区别在于，Fid用来形容一类Funcion策略的ID，
	//而KisId则为在KisFlow中KisFunction已经实例化的 实例对象ID 这个ID是随机生成且唯一
	KisId string
}

// Call
// BaseFunction 为空实现，目的为了让其他具体类型的KisFunction，如KisFunction_V 来继承BaseFuncion来重写此方法
func (base *BaseFunction) Call(ctx context.Context, flow *flow.KisFlow) error { return nil }

func (base *BaseFunction) Next() KisFunction {
	return base.N
}

func (base *BaseFunction) Prev() KisFunction {
	return base.P
}

func (base *BaseFunction) SetN(f KisFunction) {
	base.N = f
}

func (base *BaseFunction) SetP(f KisFunction) {
	base.P = f
}

func (base *BaseFunction) SetConfig(s *config.KisFuncConfig) error {
	if s == nil {
		return errors.New("KisFuncConfig is nil")
	}

	base.Config = s

	return nil
}

func (base *BaseFunction) GetId() string {
	return base.GetConfig().Fid
}

func (base *BaseFunction) GetPrevId() string {
	if base.P == nil {
		//Function为首结点
		return common.FunctionIdFirstVirtual
	}
	return base.P.GetConfig().Fid
}

func (base *BaseFunction) GetNextId() string {
	if base.N == nil {
		//Function为尾结点
		return common.FunctionIdLastVirtual
	}
	return base.N.GetConfig().Fid
}

func (base *BaseFunction) GetConfig() *config.KisFuncConfig {
	return base.Config
}

func (base *BaseFunction) SetFlow(f *flow.KisFlow) error {
	if f == nil {
		return errors.New("KisFlow is nil")
	}
	base.Flow = f
	return nil
}

func (base *BaseFunction) GetFlow() *flow.KisFlow {
	return base.Flow
}

func (base *BaseFunction) GetConnId() string {
	return base.cid
}

func (base *BaseFunction) SetConnId(id string) {
	base.cid = id
}

func (base *BaseFunction) CreateKisId() {
	base.KisId = id.KisID(common.KisIdTypeFunction)
}

func (base *BaseFunction) GetKisId() string {
	return base.KisId
}
