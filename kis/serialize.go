package kis

import (
	"kis-flow/common"
	"reflect"
)

type FaasSerialize interface {
	DecodeParam(common.KisRowArr, reflect.Type) (reflect.Value, error)
	EncodeParam(interface{}) (common.KisRowArr, error)
}
