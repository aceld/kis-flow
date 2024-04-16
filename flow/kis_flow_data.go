package flow

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/log"
	"github.com/aceld/kis-flow/metrics"
)

// CommitRow submits a single row of data to the Flow; multiple rows can be submitted multiple times
func (flow *KisFlow) CommitRow(row interface{}) error {

	flow.buffer = append(flow.buffer, row)

	return nil
}

// CommitRowBatch submits a batch of data to the Flow
func (flow *KisFlow) CommitRowBatch(rows interface{}) error {
	v := reflect.ValueOf(rows)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("Commit Data is not a slice")
	}

	for i := 0; i < v.Len(); i++ {
		row := v.Index(i).Interface().(common.KisRow)
		flow.buffer = append(flow.buffer, row)
	}

	return nil
}

// Input gets the input data for the currently executing Function in the Flow
func (flow *KisFlow) Input() common.KisRowArr {
	return flow.inPut
}

// commitSrcData submits the data source data for the current Flow, indicating the first submission of the original data source for the current Flow
// The flow's temporary data buffer is submitted to the flow's data (data is the source data backup for each Function level)
// All previous flow data will be cleared
func (flow *KisFlow) commitSrcData(ctx context.Context) error {

	// Create a batch of data
	dataCnt := len(flow.buffer)
	batch := make(common.KisRowArr, 0, dataCnt)

	for _, row := range flow.buffer {
		batch = append(batch, row)
	}

	// Clear all previous data
	flow.clearData(flow.data)

	// Record the original data for the flow for the first time
	// Because it is the first submission, PrevFunctionId is FirstVirtual because there is no upper Function
	flow.data[common.FunctionIDFirstVirtual] = batch

	// Clear the buffer
	flow.buffer = flow.buffer[0:0]

	// The first submission of data source data, for statistical total data
	if config.GlobalConfig.EnableProm == true {
		// Statistics for total data Metrics.DataTotal accumulates by 1
		metrics.Metrics.DataTotal.Add(float64(dataCnt))

		// Statistics for current Flow quantity index
		metrics.Metrics.FlowDataTotal.WithLabelValues(flow.Name).Add(float64(dataCnt))
	}

	log.Logger().DebugFX(ctx, "====> After CommitSrcData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

// getCurData gets the input data for the current Function level of the flow
func (flow *KisFlow) getCurData() (common.KisRowArr, error) {
	if flow.PrevFunctionId == "" {
		return nil, errors.New(fmt.Sprintf("flow.PrevFunctionId is not set"))
	}

	if _, ok := flow.data[flow.PrevFunctionId]; !ok {
		return nil, errors.New(fmt.Sprintf("[%s] is not in flow.data", flow.PrevFunctionId))
	}

	return flow.data[flow.PrevFunctionId], nil
}

// commitReuseData
func (flow *KisFlow) commitReuseData(ctx context.Context) error {

	// Check if there are result data from the upper layer; if not, exit the current Flow Run loop
	if len(flow.data[flow.PrevFunctionId]) == 0 {
		flow.abort = true
		return nil
	}

	// This layer's result data is equal to the upper layer's result data (reuse the upper layer's result data to this layer)
	flow.data[flow.ThisFunctionId] = flow.data[flow.PrevFunctionId]

	// Clear the buffer (If it is a ReuseData option, all the submitted data will not be carried to the next layer)
	flow.buffer = flow.buffer[0:0]

	log.Logger().DebugFX(ctx, " ====> After commitReuseData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

func (flow *KisFlow) commitVoidData(ctx context.Context) error {
	if len(flow.buffer) != 0 {
		return nil
	}

	// Create empty data
	batch := make(common.KisRowArr, 0)

	// Submit the calculated buffer data of this layer to the result data of this layer
	flow.data[flow.ThisFunctionId] = batch

	log.Logger().DebugFX(ctx, " ====> After commitVoidData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

// commitCurData submits the result data of the currently executing Function in the Flow
func (flow *KisFlow) commitCurData(ctx context.Context) error {

	// Check if this layer's calculation has result data; if not, exit the current Flow Run loop
	if len(flow.buffer) == 0 {
		flow.abort = true
		return nil
	}

	// Create a batch of data
	batch := make(common.KisRowArr, 0, len(flow.buffer))

	// If strBuf is empty, no data has been added
	for _, row := range flow.buffer {
		batch = append(batch, row)
	}

	// Submit the calculated buffer data of this layer to the result data of this layer
	flow.data[flow.ThisFunctionId] = batch

	// Clear the buffer
	flow.buffer = flow.buffer[0:0]

	log.Logger().DebugFX(ctx, " ====> After commitCurData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

// clearData clears all flow data
func (flow *KisFlow) clearData(data common.KisDataMap) {
	for k := range data {
		delete(data, k)
	}
}

func (flow *KisFlow) GetCacheData(key string) interface{} {

	if data, found := flow.cache.Get(key); found {
		return data
	}

	return nil
}

func (flow *KisFlow) SetCacheData(key string, value interface{}, Exp time.Duration) {
	if Exp == common.DefaultExpiration {
		flow.cache.Set(key, value, 0)
	} else {
		flow.cache.Set(key, value, Exp)
	}
}

// GetMetaData gets the temporary data of the current Flow object
func (flow *KisFlow) GetMetaData(key string) interface{} {
	flow.mLock.RLock()
	defer flow.mLock.RUnlock()

	data, ok := flow.metaData[key]
	if !ok {
		return nil
	}

	return data
}

// SetMetaData sets the temporary data of the current Flow object
func (flow *KisFlow) SetMetaData(key string, value interface{}) {
	flow.mLock.Lock()
	defer flow.mLock.Unlock()

	flow.metaData[key] = value
}

// GetFuncParam gets the default configuration parameters of the currently executing Function in the Flow, retrieves a key-value pair
func (flow *KisFlow) GetFuncParam(key string) string {
	flow.fplock.RLock()
	defer flow.fplock.RUnlock()

	if param, ok := flow.funcParams[flow.ThisFunctionId]; ok {
		if value, vok := param[key]; vok {
			return value
		}
	}

	return ""
}

// GetFuncParamAll gets the default configuration parameters of the currently executing Function in the Flow, retrieves all Key-Value pairs
func (flow *KisFlow) GetFuncParamAll() config.FParam {
	flow.fplock.RLock()
	defer flow.fplock.RUnlock()

	param, ok := flow.funcParams[flow.ThisFunctionId]
	if !ok {
		return nil
	}

	return param
}

// GetFuncParamsAllFuncs gets the FuncParams of all Functions in the Flow, retrieves all Key-Value pairs
func (flow *KisFlow) GetFuncParamsAllFuncs() map[string]config.FParam {
	flow.fplock.RLock()
	defer flow.fplock.RUnlock()

	return flow.funcParams
}
