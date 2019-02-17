package csv

import (
	"io"
	"testing"
)

func Test_newReader(t *testing.T) {
	type args struct {
		r       io.Reader
		options *Options
	}
	tests := []struct {
		name string
		args args
		want Options
	}{
		{
			name: "AllTheOptions",
			args: args{
				r: nil,
				options: &Options{
					Comma:            ',',
					Comment:          '#',
					LazyQuotes:       true,
					FieldsPerRecord:  2,
					TrimLeadingSpace: true,
				},
			},
			want: Options{
				Comma:            ',',
				Comment:          '#',
				LazyQuotes:       true,
				FieldsPerRecord:  2,
				TrimLeadingSpace: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newReader(tt.args.r, tt.args.options)

			if tt.want.Comma != got.Comma {
				t.Errorf("newReader() = %v, want %v", got.Comma, tt.want.Comma)
			}

			if tt.want.Comment != got.Comment {
				t.Errorf("newReader() = %v, want %v", got.Comment, tt.want.Comment)
			}

			if tt.want.LazyQuotes != got.LazyQuotes {
				t.Errorf("newReader() = %v, want %v", got.LazyQuotes, tt.want.LazyQuotes)
			}

			if tt.want.FieldsPerRecord != got.FieldsPerRecord {
				t.Errorf("newReader() = %v, want %v", got.FieldsPerRecord, tt.want.FieldsPerRecord)
			}

			if tt.want.TrimLeadingSpace != got.TrimLeadingSpace {
				t.Errorf("newReader() = %v, want %v", got.TrimLeadingSpace, tt.want.TrimLeadingSpace)
			}
		})
	}
}
