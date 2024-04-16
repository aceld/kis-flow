package common

// KisRow represents a single row of data
type KisRow interface{}

// KisRowArr represents a batch of data for a single business operation
type KisRowArr []KisRow

// KisDataMap contains all the data carried by the current Flow
// key  : Function ID where the data resides
// value: Corresponding KisRow
type KisDataMap map[string]KisRowArr
