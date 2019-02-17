package csv

import (
	"reflect"
	"strconv"
	"encoding/json"
	"encoding"
	"fmt"
)

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

var unmarshalFuncType = reflect.TypeOf(UnmarshalFunc(nil))

func nativeUnmarshalQuoted(v interface{}, data []byte) error {
	// byte -> string -> quote -> byte
	raw := []byte(strconv.Quote(string(data)))
	return json.Unmarshal(raw, v)
}

func nativeUnmarshalUnquoted(v interface{}, data []byte) error {
	return json.Unmarshal(data, v)
}

func nativeUnmarshal(t reflect.Type) UnmarshalFunc {
	if reflect.PtrTo(t).Implements(textUnmarshalerType) {
		return nativeUnmarshalQuoted;
	}

	if t.Kind() == reflect.String {
		return nativeUnmarshalQuoted;
	}

	return nativeUnmarshalUnquoted;
}

func methodUnmarshal(method reflect.Value) UnmarshalFunc {
	return func(v interface{}, data []byte) error {
		args := []reflect.Value{
			reflect.ValueOf(v),
			reflect.ValueOf(data),
		}
		responses := method.Call(args)

		err := responses[0].Interface()
		if err != nil {
			return err.(error)
		}

		return nil
	}
}



type structType struct {
	reflect.Type
}

func (s structType) getUnmarshalMethod(tags tags) (UnmarshalFunc, error) {
	if tags.Unmarshal == "" {
		return nativeUnmarshal(s.FieldByIndex(tags.index).Type), nil
	}

	// Create a zero value pointer (no reason to allocate object)
	obj := reflect.Zero(reflect.PtrTo(s.Type))
	
	// Get method on the zero object
	method := obj.MethodByName(tags.Unmarshal)
	
	
	// Verify the method existed
	if !method.IsValid() {
		return nil, fmt.Errorf("invalid method name %v", tags.Unmarshal)
	}

	// Verify signature
	if 	!(method.Type().NumIn() == 2) ||
		!(method.Type().In(0) == reflect.PtrTo(s.FieldByIndex(tags.index).Type)) ||
		!(method.Type().In(1) == reflect.TypeOf([]byte{})) ||
		!(method.Type().NumOut() == 1) ||
		!(method.Type().Out(0) == reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, fmt.Errorf("invalid method signature 1 %v", tags.Unmarshal)
	}
	
	/*
	if !method.Type().AssignableTo(unmarshalFuncType) {
		return nil, fmt.Errorf("invalid method signature %v", tags.Unmarshal)
	}
	*/
	
	return methodUnmarshal(method), nil
}