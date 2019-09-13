package csv_test

import (
	"fmt"
	"log"

	"github.com/KalleDK/go-csv/csv"
)

type CustomRecord struct {
	SimpleNumber int
	SimpleString string
	CustomMethod bool `csv:",DecodeBool"` // Telling the decoder to use the method DecodeBool to decode this bool
}

func (r *CustomRecord) DecodeBool(v *bool, text []byte) error {
	*v = "true" == string(text)
	return nil
}

var customcsv = []byte(`"SimpleNumber","SimpleString","CustomMethod"
"2","Bob","true"
"3","Sally","false"`)

func ExampleUnmarshal_custom() {
	var records []CustomRecord

	if err := csv.Unmarshal(&records, nil, customcsv); err != nil {
		log.Fatal(err)
	}

	fmt.Println(records)
	// Output:
	// [{2 Bob true} {3 Sally false}]
}
