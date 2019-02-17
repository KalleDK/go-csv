package csv_test

import (
	"fmt"
	"log"

	"github.com/KalleDK/go-csv/csv"
)

type CSVField bool



func (f *CSVField) UnmarshalText(text []byte) error {
	*f = (CSVField)(string(text) == "true")
	return nil
}

type CSVRecord struct {
	SimpleNumber   int
	SimpleString   string
	CustomMethod   bool `csv:",DecodeThird"`
	TextUnmarshal  bool
	OtherName	   string `csv:"RealName"`
}

func (r *CSVRecord) DecodeThird(v *bool, text []byte) error {
	*v = "true" == string(text)
	return nil
}

var advancedcsv = []byte(`"SimpleNumber","SimpleString","CustomMethod","TextUnmarshal","RealName"
"2","Bob","true","true","Sir Bob"
"3","Sally","false","false","Miss Alice"`)

func ExampleUnmarshal_advanced() {
	var records []CSVRecord

	if err := csv.Unmarshal(&records, nil, advancedcsv); err != nil {
		log.Fatal(err)
	}

	fmt.Println(records)
	// Output:
	// [{2 Bob true true Sir Bob} {3 Sally false false Miss Alice}]
}
