package common

import "time"

// Prefix string for generating KisId by users
const (
	KisIDTypeFlow      = "flow"   // KisId type for Flow
	KisIDTypeConnector = "conn"   // KisId type for Connector
	KisIDTypeFunction  = "func"   // KisId type for Function
	KisIDTypeGlobal    = "global" // KisId type for Global
	KisIDJoinChar      = "-"      // Joining character for KisId
)

const (
	// FunctionIDFirstVirtual is the virtual Function ID for the first node Function
	FunctionIDFirstVirtual = "FunctionIDFirstVirtual"
	// FunctionIDLastVirtual is the virtual Function ID for the last node Function
	FunctionIDLastVirtual = "FunctionIDLastVirtual"
)

// KisMode represents the mode of KisFunction
type KisMode string

const (
	// V is for Verify, which mainly performs data filtering, validation, field sorting, idempotence, etc.
	V KisMode = "Verify"

	// S is for Save, S Function will store data through KisConnector. Functions with the same Connector can logically merge.
	S KisMode = "Save"

	// L is for Load, L Function will load data through KisConnector. Functions with the same Connector can logically merge with corresponding S Function.
	L KisMode = "Load"

	// C is for Calculate, which can generate new fields, calculate new values, and perform data aggregation, analysis, etc.
	C KisMode = "Calculate"

	// E is for Expand, which serves as a custom feature Function for stream computing and is also the last Function in the current KisFlow, similar to Sink.
	E KisMode = "Expand"
)

// KisOnOff  Whether to enable the Flow
type KisOnOff int

const (
	// FlowEnable Enabled
	FlowEnable KisOnOff = 1
	// FlowDisable Disabled
	FlowDisable KisOnOff = 0
)

// KisConnType represents the type of KisConnector
type KisConnType string

const (
	// REDIS is the type of Redis
	REDIS KisConnType = "redis"
	// MYSQL is the type of MySQL
	MYSQL KisConnType = "mysql"
	// KAFKA is the type of Kafka
	KAFKA KisConnType = "kafka"
	// TIDB is the type of TiDB
	TIDB KisConnType = "tidb"
	// ES is the type of Elasticsearch
	ES KisConnType = "es"
)

// cache
const (
	// DeFaultFlowCacheCleanUp is the default cleanup time for Cache in KisFlow's Flow object Cache
	DeFaultFlowCacheCleanUp = 5 // unit: min
	// DefaultExpiration is the default time for GoCahce, permanent storage
	DefaultExpiration time.Duration = 0
)

// metrics
const (
	MetricsRoute string = "/metrics"

	LabelFlowName     string = "flow_name"
	LabelFlowID       string = "flow_id"
	LabelFunctionName string = "func_name"
	LabelFunctionMode string = "func_mode"

	CounterKisflowDataTotalName string = "kisflow_data_total"
	CounterKisflowDataTotalHelp string = "Total data volume of all KisFlow Flows"

	GamgeFlowDataTotalName string = "flow_data_total"
	GamgeFlowDataTotalHelp string = "Total data volume of each KisFlow FlowID data stream"

	GangeFlowScheCntsName string = "flow_schedule_cnts"
	GangeFlowScheCntsHelp string = "Number of times each KisFlow FlowID is scheduled"

	GangeFuncScheCntsName string = "func_schedule_cnts"
	GangeFuncScheCntsHelp string = "Number of times each KisFlow Function is scheduled"

	HistogramFunctionDurationName string = "func_run_duration"
	HistogramFunctionDurationHelp string = "Function execution time"

	HistogramFlowDurationName string = "flow_run_duration"
	HistogramFlowDurationHelp string = "Flow execution time"
)
