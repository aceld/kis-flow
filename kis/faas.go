package kis

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// FaaS Function as a Service

// Change the type definition from:
// type FaaS func(context.Context, Flow) error
// to:
// type FaaS func(context.Context, Flow, ...interface{}) error
// This allows passing data through variadic parameters of any type.
type FaaS interface{}

// FaaSDesc describes the FaaS callback computation function.
type FaaSDesc struct {
	Serialize                // Serialization implementation for the current Function's data input and output
	FnName    string         // Function name
	f         interface{}    // FaaS function
	fName     string         // Function name
	ArgsType  []reflect.Type // Function parameter types (collection)
	ArgNum    int            // Number of function parameters
	FuncType  reflect.Type   // Function type
	FuncValue reflect.Value  // Function value (function address)
}

// NewFaaSDesc creates an instance of FaaSDesc description based on the registered FnName and FaaS callback function.
func NewFaaSDesc(fnName string, f FaaS) (*FaaSDesc, error) {

	// Serialization instance
	var serializeImpl Serialize

	// Callback function value (function address)
	funcValue := reflect.ValueOf(f)

	// Callback function type
	funcType := funcValue.Type()

	// Check if the provided FaaS pointer is a function type
	if !isFuncType(funcType) {
		return nil, fmt.Errorf("provided FaaS type is %s, not a function", funcType.Name())
	}

	// Check if the FaaS function has a return value that only includes (error)
	if funcType.NumOut() != 1 || funcType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("function must have exactly one return value of type error")
	}

	// FaaS function parameter types
	argsType := make([]reflect.Type, funcType.NumIn())

	// Get the FaaS function name
	fullName := runtime.FuncForPC(funcValue.Pointer()).Name()

	// Ensure that the FaaS function parameter list contains context.Context and kis.Flow
	// Check if the function contains a parameter of type kis.Flow
	containsKisFlow := false
	// Check if the function contains a parameter of type context.Context
	containsCtx := false

	// Iterate over the FaaS function parameter types
	for i := 0; i < funcType.NumIn(); i++ {

		// Get the i-th formal parameter type
		paramType := funcType.In(i)

		if isFlowType(paramType) {
			// Check if the function contains a parameter of type kis.Flow
			containsKisFlow = true

		} else if isContextType(paramType) {
			// Check if the function contains a parameter of type context.Context
			containsCtx = true

		} else if isSliceType(paramType) {

			// Get the element type of the current parameter Slice
			itemType := paramType.Elem()

			// If the current parameter is a pointer type, get the struct type that the pointer points to
			if itemType.Kind() == reflect.Ptr {
				itemType = itemType.Elem() // Get the struct type that the pointer points to
			}

			// Check if f implements Serialize interface
			if isSerialize(itemType) {
				// If the current parameter implements the Serialize interface, use the serialization implementation of the current parameter
				serializeImpl = reflect.New(itemType).Interface().(Serialize)

			} else {
				// If the current parameter does not implement the Serialize interface, use the default serialization implementation
				serializeImpl = defaultSerialize // Use global default implementation
			}
		}

		// Append the current parameter type to the argsType collection
		argsType[i] = paramType
	}

	if !containsKisFlow {
		// If the function parameter list does not contain a parameter of type kis.Flow, return an error
		return nil, errors.New("function parameters must have kis.Flow param, please use FaaS type like: [type FaaS func(context.Context, Flow, ...interface{}) error]")
	}

	if !containsCtx {
		// If the function parameter list does not contain a parameter of type context.Context, return an error
		return nil, errors.New("function parameters must have context, please use FaaS type like: [type FaaS func(context.Context, Flow, ...interface{}) error]")
	}

	// Return the FaaSDesc description instance
	return &FaaSDesc{
		Serialize: serializeImpl,
		FnName:    fnName,
		f:         f,
		fName:     fullName,
		ArgsType:  argsType,
		ArgNum:    len(argsType),
		FuncType:  funcType,
		FuncValue: funcValue,
	}, nil
}

// isFuncType checks whether the provided paramType is a function type
func isFuncType(paramType reflect.Type) bool {
	return paramType.Kind() == reflect.Func
}

// isFlowType checks whether the provided paramType is of type kis.Flow
func isFlowType(paramType reflect.Type) bool {
	var flowInterfaceType = reflect.TypeOf((*Flow)(nil)).Elem()

	return paramType.Implements(flowInterfaceType)
}

// isContextType checks whether the provided paramType is of type context.Context
func isContextType(paramType reflect.Type) bool {
	typeName := paramType.Name()

	return strings.Contains(typeName, "Context")
}

// isSliceType checks whether the provided paramType is a slice type
func isSliceType(paramType reflect.Type) bool {
	return paramType.Kind() == reflect.Slice
}
