package csv

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

//region Simple Test

type Simple struct {
	Name string
	Age  int
}

var SimpleHeaders = []string{"Name", "Age"}

var SimpleCSV = []byte(`"Bob", 12
"Alice", 13`)

var SimpleExpected = &[]Simple{
	{"Bob", 12},
	{"Alice", 13},
}

//#endregion

//region TextUnmarshaler Test

type Age int

func (a *Age) UnmarshalText(text []byte) error {
	*a = (Age)(len(text))
	return nil
}

type TextUnmarshallerStruct struct {
	Name string
	Age  Age
}

var TextUnmarshallerHeaders = []string{"Name", "Age"}

var TextUnmarshallerCSV = []byte(`"Bob",123
"Alice",13`)

var TextUnmarshallerExpected = &[]TextUnmarshallerStruct{
	{"Bob", 3},
	{"Alice", 2},
}

//#endregion

//region CustomUnmarshal Test

type CustomUnmarshal struct {
	Name string `csv:"Name,UnmarshalName"`
	Age  int
}

func (c *CustomUnmarshal) UnmarshalName(name *string, text []byte) error {
	*name = fmt.Sprintf("Sir %v %v", string(text), len(text))
	return nil
}

/*
func (c *CustomUnmarshal) UnmarshalName(name interface{}, text []byte) error {
	*name.(*string) = fmt.Sprintf("Sir %v %v", string(text), len(text))
	return nil
}*/

var CustomUnmarshalHeaders = []string{"Name", "Age"}

var CustomUnmarshalCSV = []byte(`"Bob",12
"Alice",13`)

var CustomUnmarshalExpected = &[]CustomUnmarshal{
	{"Sir Bob 3", 12},
	{"Sir Alice 5", 13},
}

//#endregion

//region ErrorMissingRequired

type ErrorMissingRequired struct {
	Name string `csv:",,,required"`
	Age  int
}

var ErrorMissingRequiredHeaders = []string{"Age"}

//#endregion

//region ErrorMissingMethod

type ErrorMissingMethod struct {
	Name string `csv:",UnmarshalName"`
	Age  Age
}

var ErrorMissingMethodHeaders = []string{"Name"}

//#endregion

//region ErrorInvalidMethod

type ErrorInvalidMethod struct {
	Name string `csv:",UnmarshalName"`
	Age  Age
}

var ErrorInvalidMethodHeaders = []string{"Name"}

func (c *ErrorInvalidMethod) UnmarshalName(name interface{}, text string) error {
	return nil
}

//#endregion

//region ErrorSimpleUnmarshal
type ErrorSimple struct {
	Name string
	Age  int
}

var ErrorSimpleHeaders = []string{"Name", "Age"}

var ErrorSimpleCSV = []byte(`"Bob", "three"
"Alice", "four"`)

var ErrorSimpleExpected = &[]ErrorSimple{}

//#endregion

//region ErrorTextUnmarshaler Test

type ErrorAge int

func (a *ErrorAge) UnmarshalText(text []byte) error {
	return fmt.Errorf("invalid text for ErrorAge")
}

type ErrorTextUnmarshaler struct {
	Name string
	Age  ErrorAge
}

var ErrorTextUnmarshalerHeaders = []string{"Name", "Age"}

var ErrorTextUnmarshalerCSV = []byte(`"Bob",123
"Alice",13`)

var ErrorTextUnmarshalerExpected = &[]ErrorTextUnmarshaler{}

//#endregion

//region ErrorCustomUnmarshal Test

type ErrorCustomUnmarshal struct {
	Name string `csv:",UnmarshalName"`
	Age  int
}

func (c *ErrorCustomUnmarshal) UnmarshalName(name *string, text []byte) error {
	return fmt.Errorf("ErrorCustomUnmarshal Name invalid")
}

var ErrorCustomUnmarshalHeaders = []string{"Name", "Age"}

var ErrorCustomUnmarshalCSV = []byte(`"Bob",12
"Alice",13`)

var ErrorCustomUnmarshalExpected = &[]ErrorCustomUnmarshal{}

//#endregion

//region OptionalField

type OptionalField struct {
	Name string
	Age  int
}

var OptionalFieldHeaders = []string{"Age"}

var OptionalFieldCSV = []byte(`12
13`)

var OptionalFieldExpected = &[]OptionalField{
	{Age: 12},
	{Age: 13},
}

//#endregion

//region InvalidData
type InvalidData struct {
	Name string
	Age  int
}

var InvalidDataHeaders = []string{"Name", "Age"}

var InvalidDataCSV = []byte(`"Bob", 12
"Alice"`)

var InvalidDataExpected = &[]InvalidData{}

//endregion

//region InvalidDataFirst
type InvalidDataFirst struct {
	Name string
	Age  int
}

var InvalidDataFirstHeaders = []string{"Name", "Age"}

var InvalidDataFirstCSV = []byte(`"Bob"
"Alice"`)

var InvalidDataFirstExpected = &[]InvalidDataFirst{}

//endregion

func TestUnmarshal(t *testing.T) {
	type args struct {
		v       interface{}
		headers []string
		data    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Simple",
			args: args{
				v:       &[]Simple{},
				headers: SimpleHeaders,
				data:    SimpleCSV,
			},
			want:    SimpleExpected,
			wantErr: false,
		},
		{
			name: "TextUnmarshaller",
			args: args{
				v:       &[]TextUnmarshallerStruct{},
				headers: TextUnmarshallerHeaders,
				data:    TextUnmarshallerCSV,
			},
			want:    TextUnmarshallerExpected,
			wantErr: false,
		},
		{
			name: "CustomUnmarshal",
			args: args{
				v:       &[]CustomUnmarshal{},
				headers: CustomUnmarshalHeaders,
				data:    CustomUnmarshalCSV,
			},
			want:    CustomUnmarshalExpected,
			wantErr: false,
		},
		{
			name: "ErrorMissingRequired",
			args: args{
				v:       &[]ErrorMissingRequired{},
				headers: ErrorMissingRequiredHeaders,
				data:    nil,
			},
			want:    &[]ErrorMissingRequired{},
			wantErr: true,
		},
		{
			name: "ErrorMissingMethod",
			args: args{
				v:       &[]ErrorMissingMethod{},
				headers: ErrorMissingMethodHeaders,
				data:    nil,
			},
			want:    &[]ErrorMissingMethod{},
			wantErr: true,
		},
		{
			name: "ErrorInvalidMethod",
			args: args{
				v:       &[]ErrorInvalidMethod{},
				headers: ErrorInvalidMethodHeaders,
				data:    nil,
			},
			want:    &[]ErrorInvalidMethod{},
			wantErr: true,
		},
		{
			name: "ErrorTextUnmarshaler",
			args: args{
				v:       &[]ErrorTextUnmarshaler{},
				headers: ErrorTextUnmarshalerHeaders,
				data:    ErrorTextUnmarshalerCSV,
			},
			want:    &[]ErrorTextUnmarshaler{},
			wantErr: true,
		},
		{
			name: "ErrorSimpleUnmarshal",
			args: args{
				v:       &[]ErrorSimple{},
				headers: ErrorSimpleHeaders,
				data:    ErrorSimpleCSV,
			},
			want:    &[]ErrorSimple{},
			wantErr: true,
		},
		{
			name: "ErrorCustomUnmarshal",
			args: args{
				v:       &[]ErrorCustomUnmarshal{},
				headers: ErrorCustomUnmarshalHeaders,
				data:    ErrorCustomUnmarshalCSV,
			},
			want:    ErrorCustomUnmarshalExpected,
			wantErr: true,
		},
		{
			name: "OptionalField",
			args: args{
				v:       &[]OptionalField{},
				headers: OptionalFieldHeaders,
				data:    OptionalFieldCSV,
			},
			want:    OptionalFieldExpected,
			wantErr: false,
		},
		{
			name: "EmptyData",
			args: args{
				v:       &[]Simple{},
				headers: nil,
				data:    []byte{},
			},
			want:    &[]Simple{},
			wantErr: true,
		},
		{
			name: "InvalidData",
			args: args{
				v:       &[]InvalidData{},
				headers: InvalidDataHeaders,
				data:    InvalidDataCSV,
			},
			want:    InvalidDataExpected,
			wantErr: true,
		},
		{
			name: "InvalidDataFirst",
			args: args{
				v:       &[]InvalidDataFirst{},
				headers: InvalidDataFirstHeaders,
				data:    InvalidDataFirstCSV,
			},
			want:    InvalidDataFirstExpected,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.v, &Options{Headers: tt.args.headers}, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.args.v, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", tt.args.v, tt.want)
			}
		})
	}
}

type stubioreader struct {
	Value string
}

func (r *stubioreader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func TestNewDecoder(t *testing.T) {
	type args struct {
		reader  io.Reader
		headers []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Decoder
		wantErr bool
	}{
		{
			name: "InvalidReaderEOF",
			args: args{
				reader:  &stubioreader{"flaf"},
				headers: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "InvalidReaderNil",
			args: args{
				reader:  nil,
				headers: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDecoder(tt.args.reader, &Options{Headers: tt.args.headers})
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDecoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecoder() = %v, want %v", got, tt.want)
			}
		})
	}
}
