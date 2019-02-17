package csv

type fieldDecoder struct {
	recordIndex int
	structIndex []int
	unmarshal   UnmarshalFunc
}

func (d *fieldDecoder) decode(object structRecord, record csvRecord) error {
	// Field on object
	objField := object.GetField(d.structIndex)

	// Field in csv
	csvField := record[d.recordIndex]

	// Unmarshal func
	unmarshalMethod := d.unmarshal

	if err := unmarshalMethod(objField, csvField); err != nil {
		return err
	}

	return nil
}