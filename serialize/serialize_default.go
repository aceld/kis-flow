package serialize

/*
	DefaultSerialize implements the Serialize interface,
	which is used to serialize KisRowArr into a specified type, or serialize a specified type into KisRowArr.

    This section is the default serialization method provided by KisFlow, and it defaults to json serialization.

	Developers can implement their own serialization methods according to their needs.
*/

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aceld/kis-flow/common"
)

type DefaultSerialize struct{}

// UnMarshal is used to deserialize KisRowArr into a specified type.
func (f *DefaultSerialize) UnMarshal(arr common.KisRowArr, r reflect.Type) (reflect.Value, error) {
	// Ensure the passed-in type is a slice
	if r.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("r must be a slice")
	}

	slice := reflect.MakeSlice(r, 0, len(arr))

	// Iterate through each element and attempt deserialization
	for _, row := range arr {
		var elem reflect.Value
		var err error

		// Try to assert as a struct or pointer
		elem, err = unMarshalStruct(row, r.Elem())
		if err == nil {
			slice = reflect.Append(slice, elem)
			continue
		}

		// Try to directly deserialize the string
		elem, err = unMarshalJsonString(row, r.Elem())
		if err == nil {
			slice = reflect.Append(slice, elem)
			continue
		}

		// Try to serialize to JSON first and then deserialize
		elem, err = unMarshalJsonStruct(row, r.Elem())
		if err == nil {
			slice = reflect.Append(slice, elem)
		} else {
			return reflect.Value{}, fmt.Errorf("failed to decode row: %v", err)
		}
	}

	return slice, nil
}

// Try to assert as a struct or pointer
func unMarshalStruct(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
	// Check if row is of struct or struct pointer type
	rowType := reflect.TypeOf(row)
	if rowType == nil {
		return reflect.Value{}, fmt.Errorf("row is nil pointer")
	}
	if rowType.Kind() != reflect.Struct && rowType.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("row must be a struct or struct pointer type")
	}

	// If row is a pointer type, get the type it points to
	if rowType.Kind() == reflect.Ptr {
		// Nil pointer
		if reflect.ValueOf(row).IsNil() {
			return reflect.Value{}, fmt.Errorf("row is nil pointer")
		}

		// Dereference
		row = reflect.ValueOf(row).Elem().Interface()

		// Get the type after dereferencing
		rowType = reflect.TypeOf(row)
	}

	// Check if row can be asserted to elemType(target type)
	if !rowType.AssignableTo(elemType) {
		return reflect.Value{}, fmt.Errorf("row type cannot be asserted to elemType")
	}

	// Convert row to reflect.Value and return
	return reflect.ValueOf(row), nil
}

// Try to directly deserialize the string(Deserialize the Json string into a struct)
func unMarshalJsonString(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
	// Check if the source data can be asserted as a string
	str, ok := row.(string)
	if !ok {
		return reflect.Value{}, fmt.Errorf("not a string")
	}

	// Create a new struct instance to store the deserialized value
	elem := reflect.New(elemType).Elem()

	// Try to deserialize the json string into a struct.
	if err := json.Unmarshal([]byte(str), elem.Addr().Interface()); err != nil {
		return reflect.Value{}, fmt.Errorf("failed to unmarshal string to struct: %v", err)
	}

	return elem, nil
}

// Try to serialize to JSON first and then deserialize(Serialize the struct to Json string, and then deserialize the Json string into a struct)
func unMarshalJsonStruct(row common.KisRow, elemType reflect.Type) (reflect.Value, error) {
	// Serialize row to JSON string
	jsonBytes, err := json.Marshal(row)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("failed to marshal row to JSON: %v  ", err)
	}

	// Create a new struct instance to store the deserialized value
	elem := reflect.New(elemType).Interface()

	// Deserialize the JSON string into a struct
	if err := json.Unmarshal(jsonBytes, elem); err != nil {
		return reflect.Value{}, fmt.Errorf("failed to unmarshal JSON to element: %v  ", err)
	}

	return reflect.ValueOf(elem).Elem(), nil
}

// Marshal is used to serialize a specified type into KisRowArr(json serialization).
func (f *DefaultSerialize) Marshal(i interface{}) (common.KisRowArr, error) {
	var arr common.KisRowArr

	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Array:
		slice := reflect.ValueOf(i)
		for i := 0; i < slice.Len(); i++ {
			// Serialize each element to a JSON string and append it to the slice.
			jsonBytes, err := json.Marshal(slice.Index(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("failed to marshal element to JSON: %v  ", err)
			}
			arr = append(arr, string(jsonBytes))
		}
	default:
		// If it's not a slice or array type, serialize the entire struct to a JSON string directly.
		jsonBytes, err := json.Marshal(i)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal element to JSON: %v  ", err)
		}
		arr = append(arr, string(jsonBytes))
	}

	return arr, nil
}
