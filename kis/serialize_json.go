package kis

import (
	"encoding/json"
	"fmt"
	"kis-flow/common"
	"reflect"
)

type DefaultFaasSerialize struct {
}

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
			// 如果失败，则尝试直接反序列化为字符串
			elem, err = decodeString(row)
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
	elem := reflect.New(elemType).Elem()

	// 如果元素是一个结构体或指针类型，则尝试断言
	if structElem, ok := row.(reflect.Value); ok && structElem.Type().AssignableTo(elemType) {
		elem.Set(structElem)
		return elem, nil
	}

	return reflect.Value{}, fmt.Errorf("not a struct or pointer")
}

// 尝试直接反序列化字符串
func decodeString(row common.KisRow) (reflect.Value, error) {
	if str, ok := row.(string); ok {
		var intValue int
		if _, err := fmt.Sscanf(str, "%d", &intValue); err == nil {
			return reflect.ValueOf(intValue), nil
		}
	}

	return reflect.Value{}, fmt.Errorf("not a string  ")
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
