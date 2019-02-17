package csv_test

import (
	"fmt"
	"log"

	"github.com/KalleDK/go-csv/csv"
)

type SimpleRecord struct
{
	Name string
	Age int
}

var simplecsv = []byte(`"Name","Age"
"Bob","12"
"Sally","13"
"Alice","10"
`)


func ExampleUnmarshal() {
	var records []SimpleRecord
	
	if err := csv.Unmarshal(&records, nil, simplecsv); err != nil {
		log.Fatal(err)
	}

	fmt.Println(records)	

	// Output:
	// [{Bob 12} {Sally 13} {Alice 10}]
}
