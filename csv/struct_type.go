package csv

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

var unmarshalFuncType = reflect.TypeOf(UnmarshalFunc(nil))

var errorType = reflect.TypeOf((*error)(nil)).Elem()

var bytesliceType = reflect.TypeOf([]byte{})

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
		return nativeUnmarshalQuoted
	}

	if t.Kind() == reflect.String {
		return nativeUnmarshalQuoted
	}

	return nativeUnmarshalUnquoted
}

func verifyMethodSignature(methodType reflect.Type, fieldType reflect.Type) error {

	argsIn := []reflect.Type{
		reflect.PtrTo(fieldType), // First args should be a pointer to the type we want to unmarshal
		bytesliceType,            // Second args should be []byte
	}

	argsOut := []reflect.Type{
		errorType, // Only return an error
	}

	wantedMethodType := reflect.FuncOf(argsIn, argsOut, false)

	if wantedMethodType != methodType {
		return fmt.Errorf("invalid method signature %v want %v", methodType, wantedMethodType)
	}

	return nil
}

func makeUnmarshalMethod(method reflect.Value) UnmarshalFunc {
	return func(v interface{}, data []byte) error {

		// Prepare args
		args := []reflect.Value{
			reflect.ValueOf(v),
			reflect.ValueOf(data),
		}

		// Execute unmarshal
		responses := method.Call(args)

		// Forward error if any
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

func (s structType) getUnmarshalMethod(field fieldInfo) (UnmarshalFunc, error) {

	if field.Unmarshal == "" {
		return nativeUnmarshal(field.Type), nil
	}

	// Verify that the method is not a value method
	if invalidMethod, foundInvalid := s.Type.MethodByName(field.Unmarshal); foundInvalid {
		return nil, fmt.Errorf("invalid method %v can't be value method", invalidMethod)
	}

	// Create a zero value pointer (no reason to allocate object)
	obj := reflect.Zero(reflect.PtrTo(s.Type))

	// Get method on the zero object
	method := obj.MethodByName(field.Unmarshal)

	// Verify the method existed
	if !method.IsValid() {
		return nil, fmt.Errorf("invalid method name %v", field.Unmarshal)
	}

	// Verify method
	if err := verifyMethodSignature(method.Type(), field.Type); err != nil {
		return nil, err
	}

	return makeUnmarshalMethod(method), nil
}
