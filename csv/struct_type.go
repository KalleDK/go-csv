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

type objectUnmarshaler interface {
	Unmarshal(v interface{}, text []byte) error
}

type nativeUnmarshaller func(v interface{}, text []byte) error

func (n nativeUnmarshaller) Unmarshal(v interface{}, text []byte) error {
	return n(v, text)
}

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

func nativeUnmarshal(t reflect.Type) nativeUnmarshaller {
	if reflect.PtrTo(t).Implements(textUnmarshalerType) {
		return nativeUnmarshaller(nativeUnmarshalQuoted)
	}

	if t.Kind() == reflect.String {
		return nativeUnmarshaller(nativeUnmarshalQuoted)
	}

	return nativeUnmarshaller(nativeUnmarshalUnquoted)
}

func verifyMethodSignature(methodType reflect.Type, parentType reflect.Type, fieldType reflect.Type) error {

	argsIn := []reflect.Type{
		reflect.PtrTo(parentType),
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

type customUnmarshaler struct {
	obj    reflect.Value
	method reflect.Method
}

func (c customUnmarshaler) Unmarshal(v interface{}, text []byte) error {
	// Prepare args
	args := []reflect.Value{
		c.obj,
		reflect.ValueOf(v),
		reflect.ValueOf(text),
	}

	// Execute unmarshal
	responses := c.method.Func.Call(args)

	// Forward error if any
	err := responses[0].Interface()
	if err != nil {
		return err.(error)
	}

	return nil
}

type structType struct {
	reflect.Type
}

func (s structType) getUnmarshaler(field fieldInfo) (objectUnmarshaler, error) {

	if field.Unmarshal == "" {
		return nativeUnmarshal(field.Type), nil
	}

	// Verify that the method is not a value method
	if invalidMethod, foundInvalid := s.Type.MethodByName(field.Unmarshal); foundInvalid {
		return nil, fmt.Errorf("invalid method %v can't be value method", invalidMethod)
	}

	methodType, ok := reflect.PtrTo(s.Type).MethodByName(field.Unmarshal)
	if !ok {
		return nil, fmt.Errorf("invalid method name %v", field.Unmarshal)
	}

	// Verify method
	if err := verifyMethodSignature(methodType.Type, s.Type, field.Type); err != nil {
		return nil, err
	}

	// Create a zero value pointer (no reason to allocate object)
	obj := reflect.Zero(reflect.PtrTo(s.Type))

	return customUnmarshaler{
		obj:    obj,
		method: methodType,
	}, nil
}
