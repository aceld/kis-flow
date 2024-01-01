package function

import (
	"context"
	"errors"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/id"
	"kis-flow/kis"
)

type BaseFunction struct {
	Config *config.KisFuncConfig

	Flow kis.Flow //上下文环境KisFlow
	cid  string   //当前Function所依赖的KisConnectorID(如果存在)

	N kis.Function //下一个流计算Function
	P kis.Function //上一个流计算Function

	//KisId , KisFunction的实例ID，用于KisFlow内部区分不同的实例对象
	//KisId 和 Function Config中的 Fid的区别在于，Fid用来形容一类Funcion策略的ID，
	//而KisId则为在KisFlow中KisFunction已经实例化的 实例对象ID 这个ID是随机生成且唯一
	KisId string
}

// Call
// BaseFunction 为空实现，目的为了让其他具体类型的KisFunction，如KisFunction_V 来继承BaseFuncion来重写此方法
func (base *BaseFunction) Call(ctx context.Context, flow kis.Flow) error { return nil }

func (base *BaseFunction) Next() kis.Function {
	return base.N
}

func (base *BaseFunction) Prev() kis.Function {
	return base.P
}

func (base *BaseFunction) SetN(f kis.Function) {
	base.N = f
}

func (base *BaseFunction) SetP(f kis.Function) {
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

func (base *BaseFunction) SetFlow(f kis.Flow) error {
	if f == nil {
		return errors.New("KisFlow is nil")
	}
	base.Flow = f
	return nil
}

func (base *BaseFunction) GetFlow() kis.Flow {
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

// NewKisFunction 创建一个NsFunction
// flow: 当前所属的flow实例
// s : 当前function的配置策略
func NewKisFunction(flow kis.Flow, config *config.KisFuncConfig) kis.Function {
	var f kis.Function

	//工厂生产泛化对象
	switch common.KisMode(config.Fmode) {
	case common.V:
		f = new(KisFunctionV)
		break
	case common.S:
		f = new(KisFunctionS)
	case common.L:
		f = new(KisFunctionL)
	case common.C:
		f = new(KisFunctionC)
	case common.E:
		f = new(KisFunctionE)
	default:
		//LOG ERROR
		return nil
	}

	//设置基础信息属性
	if err := f.SetConfig(config); err != nil {
		panic(err)
	}

	if err := f.SetFlow(flow); err != nil {
		panic(err)
	}

	if config.Option.Cid != "" {
		f.SetConnId(config.Option.Cid)
	}

	// 生成随机实力唯一ID
	f.CreateKisId()

	return f
}
