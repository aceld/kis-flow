package function

import (
	"context"

	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/log"
)

type KisFunctionL struct {
	BaseFunction
}

func NewKisFunctionL() kis.Function {
	f := new(KisFunctionL)

	// Initialize metaData
	f.metaData = make(map[string]interface{})

	return f
}

func (f *KisFunctionL) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().DebugF("KisFunctionL, flow = %+v\n", flow)

	// Route to the specific computing Function through KisPool
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.Logger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}
