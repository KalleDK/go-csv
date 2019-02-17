package csv

import (
	"fmt"
)

type recordDecoder struct {
	decoders []*fieldDecoder
	end int
}


func newRecordDecoder(structType structType, headers headerMap) (*recordDecoder, error) {

	decoders := []*fieldDecoder{}
	last := 0

	for i := 0; i < structType.NumField(); i++ {
	
		tags := getTags(structType.Field(i))

		csvIndex, found := headers[tags.Name]
		
		if !found {
			if tags.IsOptional {
				continue	
			}
			return nil, fmt.Errorf("required field i missing in header %v", tags.Name)
		}

		if (csvIndex > last) {
			last = csvIndex
		}

		unmarshal, err := structType.getUnmarshalMethod(tags)
		if err != nil {
			return nil, err
		}

		decoders = append(
			decoders,
			&fieldDecoder{
				recordIndex: csvIndex,
				structIndex: tags.index,
				unmarshal:   unmarshal,
			},
		)

	}

	return &recordDecoder{decoders: decoders, end: last}, nil
}

func (decoder recordDecoder) Unmarshal(object structRecord, record csvRecord) error {

	if (decoder.end >= len(record)) {
		return fmt.Errorf("not enough columns in record")
	}
		
	for _, fieldDecoder := range decoder.decoders {
		if err := fieldDecoder.decode(object, record); err != nil {
			return err
		}
	}
	return nil
}
