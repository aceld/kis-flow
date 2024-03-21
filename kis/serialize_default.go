package kis

import (
	"encoding/json"
	"fmt"
	"kis-flow/common"
	"reflect"
)

type DefaultFaasSerialize struct {
}

// DecodeParam 用于将 KisRowArr 反序列化为指定类型的值。
func (f DefaultFaasSerialize) DecodeParam(arr common.KisRowArr, r reflect.Type) (reflect.Value, error) {
	// 确保传入的类型是一个切片
	if r.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("r must be a slice")
	}
	slice := reflect.MakeSlice(r, 0, len(arr))
	for _, row := range arr {
		var elem reflect.Value
		var err error

		// 先尝试断言为结构体或指针
		elem, err = decodeStruct(row, r.Elem())
		if err != nil {
			// 如果失败，则尝试直接反序列化字符串
			elem, err = decodeString(row, r.Elem())
			if err != nil {
				// 如果还失败，则尝试先序列化为 JSON 再反序列化
				elem, err = decodeJSON(row, r.Elem())
				if err != nil {
					return reflect.Value{}, fmt.Errorf("failed to decode row: %v  ", err)
				}
			}
		}

		slice = reflect.Append(slice, elem)
	}

	return slice, nil
}

// 尝试断言为结构体或指针
func decodeStruct(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
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
		if reflect.ValueOf(row).IsNil() {
			return reflect.Value{}, fmt.Errorf("row is nil pointer")
		}
		row = reflect.ValueOf(row).Elem().Interface() // 解引用
		rowType = reflect.TypeOf(row)
	}

	// 检查是否可以将 row 断言为 elemType
	if !rowType.AssignableTo(elemType) {
		return reflect.Value{}, fmt.Errorf("row type cannot be asserted to elemType")
	}

	// 将 row 转换为 reflect.Value 并返回
	return reflect.ValueOf(row), nil
}

// 尝试直接反序列化字符串
func decodeString(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
	str, ok := row.(string)
	if !ok {
		return reflect.Value{}, fmt.Errorf("not a string")
	}

	// 创建一个新的结构体实例，用于存储反序列化后的值。
	elem := reflect.New(elemType).Elem()

	// 尝试将字符串反序列化为结构体。
	if err := json.Unmarshal([]byte(str), elem.Addr().Interface()); err != nil {
		return reflect.Value{}, fmt.Errorf("failed to unmarshal string to struct: %v", err)
	}

	return elem, nil
}

// 尝试先序列化为 JSON 再反序列化
func decodeJSON(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
	jsonBytes, err := json.Marshal(row)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("failed to marshal row to JSON: %v  ", err)
	}

	elem := reflect.New(elemType).Interface()
	if err := json.Unmarshal(jsonBytes, elem); err != nil {
		return reflect.Value{}, fmt.Errorf("failed to unmarshal JSON to element: %v  ", err)
	}

	return reflect.ValueOf(elem).Elem(), nil
}

func (f DefaultFaasSerialize) EncodeParam(i interface{}) (common.KisRowArr, error) {
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
