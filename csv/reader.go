package csv

import (
	"encoding/csv"
	"io"
)

type csvReader interface {
	Read() (csvRecord, error)
}

type csvRawReader struct {
	*csv.Reader
}

func newReader(r io.Reader, options *Options) *csvRawReader {
	reader := &csvRawReader{csv.NewReader(r)}
	if options != nil {
		if options.Comma != 0 {
			reader.Comma = options.Comma
		}
		if options.Comment != 0 {
			reader.Comment = options.Comment
		}
		if options.FieldsPerRecord != 0 {
			reader.FieldsPerRecord = options.FieldsPerRecord
		}
		reader.LazyQuotes = options.LazyQuotes
		reader.TrimLeadingSpace = options.TrimLeadingSpace
	}
	return reader
}

func (r *csvRawReader) Read() (record csvRecord, err error) {
	srecord, err := r.Reader.Read()
	if err != nil {
		return nil, err
	}
	for _, field := range srecord {
		record = append(record, []byte(field))
	}

	return record, nil
}
