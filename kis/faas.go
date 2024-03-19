package kis

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
)

// FaaS Function as a Service
// type FaaS func(context.Context, *kisflow, ...interface{}) error
// 这是一个方法类型，会在注入时在方法内判断
type FaaS interface{}

type FaaSDesc struct {
	FnName    string
	f         interface{}
	fName     string
	ArgsType  []reflect.Type
	ArgNum    int
	FuncType  reflect.Type
	FuncValue reflect.Value
	FaasSerialize
}

var globalFaaSSerialize = &DefaultFaasSerialize{}

func NewFaaSDesc(fnName string, f FaaS) (*FaaSDesc, error) {
	funcValue := reflect.ValueOf(f)
	funcType := funcValue.Type()

	if err := validateFuncType(funcType, funcValue); err != nil {
		return nil, err
	}

	argsType := make([]reflect.Type, funcType.NumIn())
	fullName := runtime.FuncForPC(funcValue.Pointer()).Name()
	containsKisflowCtx := false

	for i := 0; i < funcType.NumIn(); i++ {
		paramType := funcType.In(i)
		fmt.Println(paramType.Kind(), isFlowType(paramType))
		if isFlowType(paramType) {
			containsKisflowCtx = true
		}
		argsType[i] = paramType
	}

	if !containsKisflowCtx {
		return nil, errors.New("function parameters must have Kisflow context")
	}

	if funcType.NumOut() != 1 || funcType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("function must have exactly one return value of type error")
	}

	// Check if f implements FaasSerialize interface
	var serializeImpl FaasSerialize
	if ser, ok := f.(FaasSerialize); ok {
		serializeImpl = ser
	} else {
		serializeImpl = globalFaaSSerialize // Use global default implementation
	}

	return &FaaSDesc{
		FnName:        fnName,
		f:             f,
		fName:         fullName,
		ArgsType:      argsType,
		ArgNum:        len(argsType),
		FuncType:      funcType,
		FuncValue:     funcValue,
		FaasSerialize: serializeImpl,
	}, nil
}

func validateFuncType(funcType reflect.Type, funcValue reflect.Value) error {
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("provided FaaS type is %s, not a function", funcType.Name())
	}
	return nil
}
