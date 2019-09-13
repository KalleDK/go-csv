package csv

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
)

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

var unmarshalFuncType = reflect.TypeOf(UnmarshalFunc(nil))

var errorType = reflect.TypeOf((*error)(nil)).Elem()

var bytesliceType = reflect.TypeOf([]byte{})

func nativeUnmarshalQuoted(v interface{}, data []byte) error {
	/*
		Can't get an error from a string, unless an encoder is used

		String values encode as JSON strings coerced to valid UTF-8, replacing invalid bytes with the Unicode
		replacement rune. So that the JSON will be safe to embed inside HTML <script> tags, the string is encoded using
		HTMLEscape, which replaces "<", ">", "&", U+2028, and U+2029 are escaped to "\u003c","\u003e", "\u0026",
		"\u2028", and "\u2029". This replacement can be disabled when using an Encoder, by calling
		SetEscapeHTML(false).
	*/
	raw, _ := json.Marshal(string(data))

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
