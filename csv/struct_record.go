package csv

import (
	"reflect"
)

type structField interface{}

type structRecord reflect.Value

func (r structRecord) GetField(i []int) structField {
	return reflect.Value(r).FieldByIndex(i).Addr().Interface()
}
