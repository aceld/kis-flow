package common

import "time"

// 用户生成KisId的字符串前缀
const (
	KisIdTypeFlow       = "flow"
	KisIdTypeConnnector = "conn"
	KisIdTypeFunction   = "func"
	KisIdTypeGlobal     = "global"
	KisIdJoinChar       = "-"
)

const (
	// FunctionIdFirstVirtual 为首结点Function上一层虚拟的Function ID
	FunctionIdFirstVirtual = "FunctionIdFirstVirtual"
	// FunctionIdLastVirtual 为尾结点Function下一层虚拟的Function ID
	FunctionIdLastVirtual = "FunctionIdLastVirtual"
)

type KisMode string

const (
	// V 为校验特征的KisFunction, 主要进行数据的过滤，验证，字段梳理，幂等等前置数据处理
	V KisMode = "Verify"

	// S 为存储特征的KisFunction, S会通过NsConnector进行将数据进行存储，数据的临时声明周期为NsWindow
	S KisMode = "Save"

	// L 为加载特征的KisFunction，L会通过KisConnector进行数据加载，通过该Function可以从逻辑上与对应的S Function进行并流
	L KisMode = "Load"

	// C 为计算特征的KisFunction, C会通过KisFlow中的数据计算，生成新的字段，将数据流传递给下游S进行存储，或者自己也已直接通过KisConnector进行存储
	C KisMode = "Calculate"

	// E 为扩展特征的KisFunction，作为流式计算的自定义特征Function，如，Notify 调度器触发任务的消息发送，删除一些数据，重置状态等。
	E KisMode = "Expand"
)

/*
	是否启动Flow
*/
type KisOnOff int

const (
	FlowEnable  KisOnOff = 1 // 启动
	FlowDisable KisOnOff = 0 // 不启动
)

type KisConnType string

const (
	REDIS KisConnType = "redis"
	MYSQL KisConnType = "mysql"
	KAFKA KisConnType = "kafka"
	TIDB  KisConnType = "tidb"
	ES    KisConnType = "es"
)

// cache
const (
	// DeFaultFlowCacheCleanUp KisFlow中Flow对象Cache缓存默认的清理内存时间
	DeFaultFlowCacheCleanUp = 5 //单位 min
	// DefaultExpiration 默认GoCahce时间 ，永久保存
	DefaultExpiration time.Duration = 0
)

// metrics
const (
	METRICS_ROUTE string = "/metrics"

	COUNTER_KISFLOW_DATA_TOTAL_NAME string = "kisflow_data_total"
	COUNTER_KISFLOW_DATA_TOTAL_HELP string = "KisFlow全部Flow的数据总量"
)
