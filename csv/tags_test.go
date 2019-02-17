package csv

import (
	"reflect"
	"testing"
)

func Test_tags_isOptional(t *testing.T) {
	tests := []struct {
		name string
		t    tags
		want bool
	}{
		{
			name: "AlwaysOptional",
			t:    tags{IsOptional: true},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsOptional; got != tt.want {
				t.Errorf("tags.isOptional() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTags(t *testing.T) {
	type args struct {
		field reflect.StructField
	}
	tests := []struct {
		name string
		args args
		want tags
	}{
		{
			name: "Empty",
			args: args{
				field: reflect.StructField{
					Name:  "Default",
					Index: []int{0},
					Tag:   "",
				},
			},
			want: tags{
				index:      []int{0},
				Name:       "Default",
				Unmarshal:  "",
				Marshal:    "",
				IsOptional: true,
			},
		},
		{
			name: "Full",
			args: args{
				field: reflect.StructField{
					Name:  "Default",
					Index: []int{1},
					Tag:   `csv:"MyName,MyUnmarshal,MyMarshal"`,
				},
			},
			want: tags{
				index:      []int{1},
				Name:       "MyName",
				Unmarshal:  "MyUnmarshal",
				Marshal:    "MyMarshal",
				IsOptional: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTags(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
