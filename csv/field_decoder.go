package csv

type fieldDecoder struct {
	recordIndex  int
	structIndex  []int
	unmarshaller objectUnmarshaler
}

func (d *fieldDecoder) decode(object structRecord, record csvRecord) error {
	// Field on object
	objField := object.GetField(d.structIndex)

	// Field in csv
	csvField := record[d.recordIndex]

	// Unmarshal func
	unmarshalMethod := d.unmarshaller.Unmarshal

	if err := unmarshalMethod(objField, csvField); err != nil {
		return err
	}

	return nil
}
