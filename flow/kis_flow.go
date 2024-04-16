package flow

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/conn"
	"github.com/aceld/kis-flow/function"
	"github.com/aceld/kis-flow/id"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/log"
	"github.com/aceld/kis-flow/metrics"
	"github.com/prometheus/client_golang/prometheus"

	cache "github.com/patrickmn/go-cache"
)

// KisFlow is used to manage the context environment of the entire streaming computation.
type KisFlow struct {
	// Basic information
	Id   string                // Distributed instance ID of the Flow (used internally by KisFlow to distinguish different instances)
	Name string                // Readable name of the Flow
	Conf *config.KisFlowConfig // Flow configuration policy

	// List of Functions
	Funcs          map[string]kis.Function // All managed Function objects of the current flow, key: FunctionName
	FlowHead       kis.Function            // Head of the Function list owned by the current Flow
	FlowTail       kis.Function            // Tail of the Function list owned by the current Flow
	flock          sync.RWMutex            // Lock for managing linked list insertion and reading
	ThisFunction   kis.Function            // KisFunction object currently being executed in the Flow
	ThisFunctionId string                  // ID of the Function currently being executed
	PrevFunctionId string                  // ID of the previous layer Function

	// Function list parameters
	funcParams map[string]config.FParam // Custom fixed configuration parameters of the Flow in the current Function, Key: KisID of the function instance, value: FParam
	fplock     sync.RWMutex             // Lock for managing funcParams read and write

	// Data
	buffer common.KisRowArr  // Internal buffer used to temporarily store input byte data, one data is interface{}, multiple data is []interface{} i.e. KisBatch
	data   common.KisDataMap // Data sources at various levels of the streaming computation
	inPut  common.KisRowArr  // Input data for the current Function computation
	abort  bool              // Whether to abort the Flow
	action kis.Action        // Action carried by the current Flow

	// Local cache of the flow
	cache *cache.Cache // Temporary cache context environment of the Flow

	// metaData of the flow
	metaData map[string]interface{} // Custom temporary data of the Flow
	mLock    sync.RWMutex           // Lock for managing metaData read and write
}

// NewKisFlow creates a KisFlow.
func NewKisFlow(conf *config.KisFlowConfig) kis.Flow {
	flow := new(KisFlow)
	// Instance Id
	flow.Id = id.KisID(common.KisIDTypeFlow)

	// Basic information
	flow.Name = conf.FlowName
	flow.Conf = conf

	// List of Functions
	flow.Funcs = make(map[string]kis.Function)
	flow.funcParams = make(map[string]config.FParam)

	// Data
	flow.data = make(common.KisDataMap)

	// Initialize local cache
	flow.cache = cache.New(cache.NoExpiration, common.DeFaultFlowCacheCleanUp*time.Minute)

	// Initialize temporary data
	flow.metaData = make(map[string]interface{})

	return flow
}

// Fork gets a copy (deep copy) of the Flow.
func (flow *KisFlow) Fork(ctx context.Context) kis.Flow {

	cfg := flow.Conf

	// Generate a new Flow based on the previous configuration
	newFlow := NewKisFlow(cfg)

	for _, fp := range flow.Conf.Flows {
		if _, ok := flow.funcParams[flow.Funcs[fp.FuncName].GetID()]; !ok {
			// The current function has no Params configured
			_ = newFlow.AppendNewFunction(flow.Funcs[fp.FuncName].GetConfig(), nil)
		} else {
			// The current function has configured Params
			_ = newFlow.AppendNewFunction(flow.Funcs[fp.FuncName].GetConfig(), fp.Params)
		}
	}

	log.Logger().DebugFX(ctx, "=====>Flow Fork, oldFlow.funcParams = %+v\n", flow.funcParams)
	log.Logger().DebugFX(ctx, "=====>Flow Fork, newFlow.funcParams = %+v\n", newFlow.GetFuncParamsAllFuncs())

	return newFlow
}

// Link links the Function to the Flow, and also adds the Function's configuration parameters to the Flow's configuration.
// fConf: Current Function strategy
// fParams: Dynamic parameters carried by the current Flow's Function
func (flow *KisFlow) Link(fConf *config.KisFuncConfig, fParams config.FParam) error {

	// Add Function to Flow
	_ = flow.AppendNewFunction(fConf, fParams)

	// Add Function to FlowConfig
	flowFuncParam := config.KisFlowFunctionParam{
		FuncName: fConf.FName,
		Params:   fParams,
	}
	flow.Conf.AppendFunctionConfig(flowFuncParam)

	return nil
}

// AppendNewFunction appends a new Function to the Flow.
func (flow *KisFlow) AppendNewFunction(fConf *config.KisFuncConfig, fParams config.FParam) error {
	// Create Function instance
	f := function.NewKisFunction(flow, fConf)

	if fConf.Option.CName != "" {
		// The current Function has a Connector association and needs to initialize the Connector instance

		// Get Connector configuration
		connConfig, err := fConf.GetConnConfig()
		if err != nil {
			panic(err)
		}

		// Create Connector object
		connector := conn.NewKisConnector(connConfig)

		// Initialize Connector, execute the Connector Init method
		if err = connector.Init(); err != nil {
			panic(err)
		}

		// Associate the Function instance with the Connector instance
		_ = f.AddConnector(connector)
	}

	// Add Function to Flow
	if err := flow.appendFunc(f, fParams); err != nil {
		return err
	}

	return nil
}

// appendFunc adds the Function to the Flow, linked list operation
func (flow *KisFlow) appendFunc(function kis.Function, fParam config.FParam) error {

	if function == nil {
		return errors.New("AppendFunc append nil to List")
	}

	flow.flock.Lock()
	defer flow.flock.Unlock()

	if flow.FlowHead == nil {
		// First time adding a node
		flow.FlowHead = function
		flow.FlowTail = function

		function.SetN(nil)
		function.SetP(nil)

	} else {
		// Insert the function at the end of the linked list
		function.SetP(flow.FlowTail)
		function.SetN(nil)

		flow.FlowTail.SetN(function)
		flow.FlowTail = function
	}

	// Add the detailed Function Name-Hash correspondence to the flow object
	flow.Funcs[function.GetConfig().FName] = function

	// First add the Params parameters carried by the function by default
	params := make(config.FParam)
	for key, value := range function.GetConfig().Option.Params {
		params[key] = value
	}

	// Then add the function definition parameters carried by the flow (overwriting duplicates)
	for key, value := range fParam {
		params[key] = value
	}

	// Store the obtained FParams in the flow structure for direct access by the function
	// The key is the KisId of the current Function, not using Fid to prevent adding two Functions with the same strategy Id to a Flow
	flow.funcParams[function.GetID()] = params

	return nil
}

// Run starts the streaming computation of KisFlow, starting from the initial Function.
func (flow *KisFlow) Run(ctx context.Context) error {

	var fn kis.Function

	fn = flow.FlowHead
	flow.abort = false

	if flow.Conf.Status == int(common.FlowDisable) {
		// Flow is configured to be disabled
		return nil
	}

	// Metrics
	var funcStart time.Time
	var flowStart time.Time

	// Since no Function has been executed at this time, PrevFunctionId is FirstVirtual because there is no previous layer Function
	flow.PrevFunctionId = common.FunctionIDFirstVirtual

	// Commit the original data stream
	if err := flow.commitSrcData(ctx); err != nil {
		return err
	}

	// Metrics
	if config.GlobalConfig.EnableProm == true {
		// Count the number of Flow schedules
		metrics.Metrics.FlowScheduleCntsToTal.WithLabelValues(flow.Name).Inc()
		// Count the execution time of Flow
		flowStart = time.Now()
	}

	// Streaming chain call
	for fn != nil && flow.abort == false {

		// Record the current Function being executed by the flow
		fid := fn.GetID()
		flow.ThisFunction = fn
		flow.ThisFunctionId = fid

		fName := fn.GetConfig().FName
		fMode := fn.GetConfig().FMode

		if config.GlobalConfig.EnableProm == true {
			// Count the number of Function schedules
			metrics.Metrics.FuncScheduleCntsTotal.WithLabelValues(fName, fMode).Inc()

			// Count the time consumed by Function, record the start time
			funcStart = time.Now()
		}

		// Get the source data that the current Function needs to process
		if inputData, err := flow.getCurData(); err != nil {
			log.Logger().ErrorFX(ctx, "flow.Run(): getCurData err = %s\n", err.Error())
			return err
		} else {
			flow.inPut = inputData
		}

		if err := fn.Call(ctx, flow); err != nil {
			// Error
			return err
		} else {
			// Success
			fn, err = flow.dealAction(ctx, fn)
			if err != nil {
				return err
			}

			// Count the time consumed by Function
			if config.GlobalConfig.EnableProm == true {
				// Function consumption time
				duration := time.Since(funcStart)

				// Count the current Function metrics, do time statistics
				metrics.Metrics.FunctionDuration.With(
					prometheus.Labels{
						common.LabelFunctionName: fName,
						common.LabelFunctionMode: fMode}).Observe(duration.Seconds() * 1000)
			}

		}
	}

	// Metrics
	if config.GlobalConfig.EnableProm == true {
		// Count the execution time of Flow
		duration := time.Since(flowStart)
		metrics.Metrics.FlowDuration.WithLabelValues(flow.Name).Observe(duration.Seconds() * 1000)
	}

	return nil
}

// Next the current Flow enters the Action action carried by the next layer Function.
func (flow *KisFlow) Next(acts ...kis.ActionFunc) error {

	// Load the Action actions carried by Function FaaS
	flow.action = kis.LoadActions(acts)

	return nil
}

func (flow *KisFlow) GetName() string {
	return flow.Name
}

func (flow *KisFlow) GetID() string {
	return flow.Id
}

func (flow *KisFlow) GetThisFunction() kis.Function {
	return flow.ThisFunction
}

func (flow *KisFlow) GetThisFuncConf() *config.KisFuncConfig {
	return flow.ThisFunction.GetConfig()
}

// GetConnector gets the Connector of the Function currently being executed by the Flow.
func (flow *KisFlow) GetConnector() (kis.Connector, error) {
	if connector := flow.ThisFunction.GetConnector(); connector != nil {
		return connector, nil
	} else {
		return nil, errors.New("GetConnector(): Connector is nil")
	}
}

// GetConnConf gets the Connector configuration of the Function currently being executed by the Flow.
func (flow *KisFlow) GetConnConf() (*config.KisConnConfig, error) {
	if connector := flow.ThisFunction.GetConnector(); connector != nil {
		return connector.GetConfig(), nil
	} else {
		return nil, errors.New("GetConnConf(): Connector is nil")
	}
}

func (flow *KisFlow) GetConfig() *config.KisFlowConfig {
	return flow.Conf
}

// GetFuncConfigByName gets the configuration of the current Flow by Function name.
func (flow *KisFlow) GetFuncConfigByName(funcName string) *config.KisFuncConfig {
	if f, ok := flow.Funcs[funcName]; ok {
		return f.GetConfig()
	} else {
		log.Logger().ErrorF("GetFuncConfigByName(): Function %s not found", funcName)
		return nil
	}
}
