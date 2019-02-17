package csv

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_headerList_ToMap(t *testing.T) {
	tests := []struct {
		name    string
		headers headerList
		want    headerMap
	}{
		{
			name:    "Simple Conversion",
			headers: []string{"one", "two", "three"},
			want: map[string]int{
				"one":   0,
				"two":   1,
				"three": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.headers.ToMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("headerList.ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

type stubreader []string

func (r stubreader) Read() (csvRecord, error) {
	result := csvRecord{}
	for _, field := range r {
		result = append(result, []byte(field))
	}
	return result, nil
}

type errorreader struct {
	error
}

func (r errorreader) Read() (csvRecord, error) {
	return nil, r.error
}

func Test_getHeaders(t *testing.T) {
	type args struct {
		reader  csvReader
		headers headerList
	}
	tests := []struct {
		name    string
		args    args
		want    headerMap
		wantErr bool
	}{
		{
			name: "HeaderList",
			args: args{
				reader:  stubreader{},
				headers: []string{"one", "two", "three"},
			},
			want: map[string]int{
				"one":   0,
				"two":   1,
				"three": 2,
			},
			wantErr: false,
		},
		{
			name: "HeaderFromReader",
			args: args{
				reader:  stubreader{"one", "two", "three"},
				headers: nil,
			},
			want: map[string]int{
				"one":   0,
				"two":   1,
				"three": 2,
			},
			wantErr: false,
		},
		{
			name: "HeaderError",
			args: args{
				reader:  errorreader{fmt.Errorf("flaf")},
				headers: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ReaderCantBeNil",
			args: args{
				reader:  nil,
				headers: []string{"one", "two", "three"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getHeaders(tt.args.reader, tt.args.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
