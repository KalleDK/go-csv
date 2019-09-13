package csv

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

/*
Options for how the CSV file is encoded
*/
type Options struct {
	// Headers is for mapping the Struct's fieldnames to columns
	// If Headers are nil, the first record is expected to be headers
	Headers []string

	// Comma is the field delimiter.
	// It is set to comma (',') by NewDecoder.
	// Comma must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	Comma rune

	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the
	// field, even if TrimLeadingSpace is true.
	// Comment must be a valid rune and must not be \r, \n,
	// or the Unicode replacement character (0xFFFD).
	// It must also not be equal to Comma.
	Comment rune

	// FieldsPerRecord is the number of expected fields per record.
	// If FieldsPerRecord is positive, Read requires each record to
	// have the given number of fields. If FieldsPerRecord is 0, Read sets it to
	// the number of fields in the first record, so that future records must
	// have the same field count. If FieldsPerRecord is negative, no check is
	// made and records may have a variable number of fields.
	FieldsPerRecord int

	// If LazyQuotes is true, a quote may appear in an unquoted field and a
	// non-doubled quote may appear in a quoted field.
	LazyQuotes bool

	// If TrimLeadingSpace is true, leading white space in a field is ignored.
	// This is done even if the field delimiter, Comma, is white space.
	TrimLeadingSpace bool
}

/*
A Decoder reads and decodes CSV values from an input stream.
*/
type Decoder struct {
	headers headerMap
	reader  csvReader
}

/*
Decode reads the next CSV-encoded value from its input and stores it in the value pointed to by v.

See the documentation for Unmarshal for details about the conversion of CSV into a Go value.
*/
func (d Decoder) Decode(v interface{}) error {
	valueSlice := reflect.ValueOf(v).Elem()

	decoder, err := newRecordDecoder(structType{Type: valueSlice.Type().Elem()}, d.headers)
	if err != nil {
		return err
	}

	slice := reflect.MakeSlice(valueSlice.Type(), 0, 0)
	record, err := d.reader.Read()
	for err == nil {
		value := reflect.New(valueSlice.Type().Elem()).Elem()
		err = decoder.Unmarshal(structRecord(value), record)
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, value)
		record, err = d.reader.Read()
	}

	if err != io.EOF {
		return err
	}

	valueSlice.Set(slice)

	return nil
}

/*
NewDecoder returns a new decoder that reads from r.

The decoder introduces its own buffering and may read data from r beyond the CSV values requested.

If headers is nil the headers are expected to be in the first csv record
*/
func NewDecoder(r io.Reader, options *Options) (*Decoder, error) {
	if r == nil {
		return nil, fmt.Errorf("reader can't be nil")
	}

	csvreader := newReader(r, options)

	var headerlist []string
	if options != nil {
		headerlist = options.Headers
	}

	headermap, err := getHeaders(csvreader, headerlist)
	if err != nil {
		return nil, err
	}

	return &Decoder{headermap, csvreader}, nil
}

/*
Unmarshal parses the CSV-encoded data and stores the result in the value pointed to by v. If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
*/
func Unmarshal(v interface{}, options *Options, data []byte) error {
	ioreader := strings.NewReader(string(data))

	decoder, err := NewDecoder(ioreader, options)
	if err != nil {
		return err
	}

	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}
