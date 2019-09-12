package csv

import (
	"reflect"
	"testing"
)

var intType = reflect.TypeOf(int(0))

var stringType = reflect.TypeOf(string(""))

func Test_getFieldInfo(t *testing.T) {
	type args struct {
		field reflect.StructField
	}
	tests := []struct {
		name string
		args args
		want fieldInfo
	}{
		{
			name: "Empty",
			args: args{
				field: reflect.StructField{
					Type:  intType,
					Name:  "Default",
					Index: []int{0},
					Tag:   "",
				},
			},
			want: fieldInfo{
				index:      []int{0},
				Name:       "Default",
				Unmarshal:  "",
				Marshal:    "",
				IsOptional: true,
				Type:       intType,
			},
		},
		{
			name: "Full",
			args: args{
				field: reflect.StructField{
					Type:  intType,
					Name:  "Default",
					Index: []int{1},
					Tag:   `csv:"MyName,MyUnmarshal,MyMarshal"`,
				},
			},
			want: fieldInfo{
				index:      []int{1},
				Name:       "MyName",
				Unmarshal:  "MyUnmarshal",
				Marshal:    "MyMarshal",
				IsOptional: true,
				Type:       intType,
			},
		},
		{
			name: "Required",
			args: args{
				field: reflect.StructField{
					Type:  stringType,
					Name:  "Default",
					Index: []int{1},
					Tag:   `csv:"MyName,MyUnmarshal,MyMarshal,required"`,
				},
			},
			want: fieldInfo{
				index:      []int{1},
				Name:       "MyName",
				Unmarshal:  "MyUnmarshal",
				Marshal:    "MyMarshal",
				IsOptional: false,
				Type:       stringType,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldInfo(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFieldInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
