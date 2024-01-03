package common

// KisRow 一行数据
type KisRow interface{}

// KisRowArr 一次业务的批量数据
type KisRowArr []KisRow

/*
	KisDataMap 当前Flow承载的全部数据
   	key	:  数据所在的Function ID
    value: 对应的KisRow
*/
type KisDataMap map[string]KisRowArr
