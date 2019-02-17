package csv

import (
	"fmt"
	"reflect"
	"testing"
)

type ValidObject struct {
	FirstField  int
	SecondField int
}

func newValuePtr(v interface{}) reflect.Value {
	return reflect.ValueOf(v)
}

func newValue(v interface{}) structRecord {
	return structRecord(newValuePtr(v).Elem())
}

func unmarshalStubCreator(validObject interface{}) func(int) UnmarshalFunc {
	return func(idx int) UnmarshalFunc {
		validField := []byte(validRecord[idx])
		return func(v interface{}, data []byte) error {
			validPtr := reflect.ValueOf(validObject).Elem().UnsafeAddr()
			vPtr := reflect.ValueOf(v).Elem().UnsafeAddr()
			if validPtr != vPtr {
				return fmt.Errorf("invalid object %v %v", validPtr, vPtr)
			}

			if !reflect.DeepEqual(data, validField) {
				return fmt.Errorf("invalid record %v", data)
			}

			return nil
		}
	}

}

var validObject = ValidObject{}
var invalidObject = ValidObject{}
var validRecord = [][]byte{[]byte(`"first"`), []byte(`"second"`)}
var invalidRecord = [][]byte{[]byte(`"third"`), []byte(`"fourth"`)}
var unmarshalFirstFieldStub = unmarshalStubCreator(&validObject.FirstField)
var unmarshalSecondFieldStub = unmarshalStubCreator(&validObject.SecondField)

func Test_unmarshal(t *testing.T) {
	type args struct {
		decoders []*fieldDecoder
		object   structRecord
		record   [][]byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid",
			args: args{
				decoders: []*fieldDecoder{
					&fieldDecoder{
						recordIndex: 0,
						structIndex: []int{0},
						unmarshal:   unmarshalFirstFieldStub(0),
					},
					&fieldDecoder{
						recordIndex: 1,
						structIndex: []int{1},
						unmarshal:   unmarshalSecondFieldStub(1),
					},
				},
				object: newValue(&validObject),
				record: validRecord,
			},
			wantErr: false,
		},
		{
			name: "Invalid",
			args: args{
				decoders: []*fieldDecoder{
					&fieldDecoder{
						recordIndex: 0,
						structIndex: []int{0},
						unmarshal:   unmarshalFirstFieldStub(1),
					},
					&fieldDecoder{
						recordIndex: 1,
						structIndex: []int{1},
						unmarshal:   unmarshalSecondFieldStub(1),
					},
				},
				object: newValue(&validObject),
				record: validRecord,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := recordDecoder{tt.args.decoders, 1}
			if err := dec.Unmarshal(tt.args.object, tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStructDecoder_Unmarshal(t *testing.T) {
	type args struct {
		object structRecord
		record [][]byte
	}
	tests := []struct {
		name    string
		decoder recordDecoder
		args    args
		wantErr bool
	}{
		{
			name: "Valid",
			decoder: recordDecoder{
				decoders: []*fieldDecoder{
					&fieldDecoder{
						recordIndex: 0,
						structIndex: []int{0},
						unmarshal:   unmarshalFirstFieldStub(0),
					},
					&fieldDecoder{
						recordIndex: 1,
						structIndex: []int{1},
						unmarshal:   unmarshalSecondFieldStub(1),
					},
				},
			},
			args: args{
				object: newValue(&validObject),
				record: validRecord,
			},
			wantErr: false,
		},
		{
			name: "Invalid",
			decoder: recordDecoder{
				decoders: []*fieldDecoder{
					&fieldDecoder{
						recordIndex: 0,
						structIndex: []int{0},
						unmarshal:   unmarshalFirstFieldStub(0),
					},
					&fieldDecoder{
						recordIndex: 1,
						structIndex: []int{1},
						unmarshal:   unmarshalSecondFieldStub(0),
					},
				},
			},
			args: args{
				object: newValue(&validObject),
				record: validRecord,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.decoder.Unmarshal(tt.args.object, tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("StructDecoder.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
