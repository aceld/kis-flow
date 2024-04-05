/*
		DefaultSerialize 实现了 Serialize 接口，用于将 KisRowArr 序列化为指定类型的值，或将指定类型的值序列化为 KisRowArr。
	    这部分是KisFlow默认提供的序列化办法，默认均是josn序列化，开发者可以根据自己的需求实现自己的序列化办法。
*/
package serialize

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aceld/kis-flow/common"
)

type DefaultSerialize struct{}

// UnMarshal 用于将 KisRowArr 反序列化为指定类型的值。
func (f *DefaultSerialize) UnMarshal(arr common.KisRowArr, r reflect.Type) (reflect.Value, error) {
	// 确保传入的类型是一个切片
	if r.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("r must be a slice")
	}

	slice := reflect.MakeSlice(r, 0, len(arr))

	// 遍历每个元素并尝试反序列化
	for _, row := range arr {
		var elem reflect.Value
		var err error

		// 尝试断言为结构体或指针
		elem, err = unMarshalStruct(row, r.Elem())
		if err == nil {
			slice = reflect.Append(slice, elem)
			continue
		}

		// 尝试直接反序列化字符串
		elem, err = unMarshalJsonString(row, r.Elem())
		if err == nil {
			slice = reflect.Append(slice, elem)
			continue
		}

		// 尝试先序列化为 JSON 再反序列化
		elem, err = unMarshalJsonStruct(row, r.Elem())
		if err == nil {
			slice = reflect.Append(slice, elem)
		} else {
			return reflect.Value{}, fmt.Errorf("failed to decode row: %v", err)
		}
	}

	return slice, nil
}

// 尝试断言为结构体或指针
func unMarshalStruct(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
	// 检查 row 是否为结构体或结构体指针类型
	rowType := reflect.TypeOf(row)
	if rowType == nil {
		return reflect.Value{}, fmt.Errorf("row is nil pointer")
	}
	if rowType.Kind() != reflect.Struct && rowType.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("row must be a struct or struct pointer type")
	}

	// 如果 row 是指针类型，则获取它指向的类型
	if rowType.Kind() == reflect.Ptr {
		// 空指针
		if reflect.ValueOf(row).IsNil() {
			return reflect.Value{}, fmt.Errorf("row is nil pointer")
		}

		// 解引用
		row = reflect.ValueOf(row).Elem().Interface()

		// 拿到解引用后的类型
		rowType = reflect.TypeOf(row)
	}

	// 检查是否可以将 row 断言为 elemType(目标类型)
	if !rowType.AssignableTo(elemType) {
		return reflect.Value{}, fmt.Errorf("row type cannot be asserted to elemType")
	}

	// 将 row 转换为 reflect.Value 并返回
	return reflect.ValueOf(row), nil
}

// 尝试直接反序列化字符串(将Json字符串 反序列化为 结构体)
func unMarshalJsonString(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
	// 判断源数据是否可以断言成string
	str, ok := row.(string)
	if !ok {
		return reflect.Value{}, fmt.Errorf("not a string")
	}

	// 创建一个新的结构体实例，用于存储反序列化后的值
	elem := reflect.New(elemType).Elem()

	// 尝试将json字符串反序列化为结构体。
	if err := json.Unmarshal([]byte(str), elem.Addr().Interface()); err != nil {
		return reflect.Value{}, fmt.Errorf("failed to unmarshal string to struct: %v", err)
	}

	return elem, nil
}

// 尝试先序列化为 JSON 再反序列化(将结构体转换成Json字符串，再将Json字符串 反序列化为 结构体)
func unMarshalJsonStruct(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
	// 将 row 序列化为 JSON 字符串
	jsonBytes, err := json.Marshal(row)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("failed to marshal row to JSON: %v  ", err)
	}

	// 创建一个新的结构体实例，用于存储反序列化后的值
	elem := reflect.New(elemType).Interface()

	// 将 JSON 字符串反序列化为结构体
	if err := json.Unmarshal(jsonBytes, elem); err != nil {
		return reflect.Value{}, fmt.Errorf("failed to unmarshal JSON to element: %v  ", err)
	}

	return reflect.ValueOf(elem).Elem(), nil
}

// Marshal 用于将指定类型的值序列化为 KisRowArr(json 序列化)。
func (f *DefaultSerialize) Marshal(i interface{}) (common.KisRowArr, error) {
	var arr common.KisRowArr

	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array:
		slice := reflect.ValueOf(i)
		for i := 0; i < slice.Len(); i++ {
			// 序列化每个元素为 JSON 字符串，并将其添加到切片中。
			jsonBytes, err := json.Marshal(slice.Index(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal element to JSON: %v  ", err)
			}
			arr = append(arr, string(jsonBytes))
		}
	default:
		// 如果不是切片或数组类型，则直接序列化整个结构体为 JSON 字符串。
		jsonBytes, err := json.Marshal(i)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal element to JSON: %v  ", err)
		}
		arr = append(arr, string(jsonBytes))
	}

	return arr, nil
}
