package kis

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// FaaS Function as a Service

// 将
// type FaaS func(context.Context, Flow) error
// 改为
// type FaaS func(context.Context, Flow, ...interface{}) error
// 可以通过可变参数的任意输入类型进行数据传递
type FaaS interface{}

// FaaSDesc FaaS 回调计算业务函数 描述
type FaaSDesc struct {
	Serialize                // 当前Function的数据输入输出序列化实现
	FnName    string         // Function名称
	f         interface{}    // FaaS 函数
	fName     string         // 函数名称
	ArgsType  []reflect.Type // 函数参数类型（集合）
	ArgNum    int            // 函数参数个数
	FuncType  reflect.Type   // 函数类型
	FuncValue reflect.Value  // 函数值(函数地址)
}

// NewFaaSDesc 根据用户注册的FnName 和FaaS 回调函数，创建 FaaSDesc 描述实例
func NewFaaSDesc(fnName string, f FaaS) (*FaaSDesc, error) {

	// 输入输出序列化实例
	var serializeImpl Serialize

	// 传入的回调函数FaaS,函数值(函数地址)
	funcValue := reflect.ValueOf(f)

	// 传入的回调函数FaaS 类型
	funcType := funcValue.Type()

	// 判断传递的FaaS指针是否是函数类型
	if !isFuncType(funcType) {
		return nil, fmt.Errorf("provided FaaS type is %s, not a function", funcType.Name())
	}

	// 判断传递的FaaS函数是否有返回值类型是只包括(error)
	if funcType.NumOut() != 1 || funcType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("function must have exactly one return value of type error")
	}

	// FaaS函数的参数类型集合
	argsType := make([]reflect.Type, funcType.NumIn())

	// 获取FaaS的函数名称
	fullName := runtime.FuncForPC(funcValue.Pointer()).Name()

	// 确保  FaaS func(context.Context, Flow, ...interface{}) error 形参列表，存在context.Context 和 kis.Flow

	// 是否包含kis.Flow类型的形参
	containsKisFlow := false
	// 是否包含context.Context类型的形参
	containsCtx := false

	// 遍历FaaS的形参类型
	for i := 0; i < funcType.NumIn(); i++ {

		// 取出第i个形式参数类型
		paramType := funcType.In(i)

		if isFlowType(paramType) {
			// 判断是否包含kis.Flow类型的形参
			containsKisFlow = true

		} else if isContextType(paramType) {
			// 判断是否包含context.Context类型的形参
			containsCtx = true

		} else if isSliceType(paramType) {

			// 获取当前参数Slice的元素类型
			itemType := paramType.Elem()

			// 如果当前参数是一个指针类型，则获取指针指向的结构体类型
			if itemType.Kind() == reflect.Ptr {
				itemType = itemType.Elem() // 获取指针指向的结构体类型
			}

			// Check if f implements Serialize interface
			// (检测传递的FaaS函数是否实现了Serialize接口)
			if isSerialize(itemType) {
				// 如果当前形参实现了Serialize接口，则使用当前形参的序列化实现
				serializeImpl = reflect.New(itemType).Interface().(Serialize)

			} else {
				// 如果当前形参没有实现Serialize接口，则使用默认的序列化实现
				serializeImpl = defaultSerialize // Use global default implementation
			}
		} else {
			// Other types are not supported
		}

		// 将当前形参类型追加到argsType集合中
		argsType[i] = paramType
	}

	if !containsKisFlow {
		// 不包含kis.Flow类型的形参，返回错误
		return nil, errors.New("function parameters must have kis.Flow param, please use FaaS type like: [type FaaS func(context.Context, Flow, ...interface{}) error]")
	}

	if !containsCtx {
		// 不包含context.Context类型的形参，返回错误
		return nil, errors.New("function parameters must have context, please use FaaS type like: [type FaaS func(context.Context, Flow, ...interface{}) error]")
	}

	// 返回FaaSDesc描述实例
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

// isFuncType 判断传递进来的 paramType 是否是函数类型
func isFuncType(paramType reflect.Type) bool {
	return paramType.Kind() == reflect.Func
}

// isFlowType 判断传递进来的 paramType 是否是 kis.Flow 类型
func isFlowType(paramType reflect.Type) bool {
	var flowInterfaceType = reflect.TypeOf((*Flow)(nil)).Elem()

	return paramType.Implements(flowInterfaceType)
}

// isContextType 判断传递进来的 paramType 是否是 context.Context 类型
func isContextType(paramType reflect.Type) bool {
	typeName := paramType.Name()

	return strings.Contains(typeName, "Context")
}

// isSliceType 判断传递进来的 paramType 是否是切片类型
func isSliceType(paramType reflect.Type) bool {
	return paramType.Kind() == reflect.Slice
}
