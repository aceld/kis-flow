package kis

import (
	"reflect"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/serialize"
)

// Serialize Data serialization interface
type Serialize interface {
	// UnMarshal is used to deserialize KisRowArr to a value of the specified type.
	UnMarshal(common.KisRowArr, reflect.Type) (reflect.Value, error)
	// Marshal is used to serialize a value of the specified type to KisRowArr.
	Marshal(interface{}) (common.KisRowArr, error)
}

// defaultSerialize Default serialization implementation provided by KisFlow (developers can customize)
var defaultSerialize = &serialize.DefaultSerialize{}

// isSerialize checks if the provided paramType implements the Serialize interface
func isSerialize(paramType reflect.Type) bool {
	return paramType.Implements(reflect.TypeOf((*Serialize)(nil)).Elem())
}
