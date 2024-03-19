:package kis

import (
	"reflect"
	"strings"
)

func isContextType(paramType reflect.Type) bool {
	typeName := paramType.Name()
	return strings.Contains(typeName, "Context")
}
