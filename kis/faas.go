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

	var serializeImpl FaasSerialize

	funcValue := reflect.ValueOf(f)
	funcType := funcValue.Type()

	if err := validateFuncType(funcType, funcValue); err != nil {
		return nil, err
	}

	argsType := make([]reflect.Type, funcType.NumIn())
	fullName := runtime.FuncForPC(funcValue.Pointer()).Name()
	containsKisFlow := false
	containsCtx := false

	for i := 0; i < funcType.NumIn(); i++ {
		paramType := funcType.In(i)
		if isFlowType(paramType) {
			containsKisFlow = true
		} else if isContextType(paramType) {
			containsCtx = true
		} else {
			itemType := paramType.Elem()
			// 如果切片元素是指针类型，则获取指针所指向的类型
			if itemType.Kind() == reflect.Ptr {
				itemType = itemType.Elem()
			}
			// Check if f implements FaasSerialize interface
			if isFaasSerialize(itemType) {
				serializeImpl = reflect.New(itemType).Interface().(FaasSerialize)
			} else {
				serializeImpl = globalFaaSSerialize // Use global default implementation
			}

		}
		argsType[i] = paramType
	}

	if !containsKisFlow {
		return nil, errors.New("function parameters must have Kisflow context")
	}
	if !containsCtx {
		return nil, errors.New("function parameters must have context")
	}

	if funcType.NumOut() != 1 || funcType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("function must have exactly one return value of type error")
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
