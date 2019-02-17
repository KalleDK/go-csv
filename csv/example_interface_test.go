package csv_test

import (
	"fmt"
	"log"

	"github.com/KalleDK/go-csv/csv"
)

type Record struct
{
	Name string
	Age int
}

var examplecsv = []byte(`"Name","Age"
"Bob","12"
"Sally","13"
"Alice","10"
`)

func Example() {
	var records []Record
	
	if err := csv.Unmarshal(&records, nil, examplecsv); err != nil {
		log.Fatal(err)
	}

	fmt.Println(records)
	// Output:
	// [{Bob 12} {Sally 13} {Alice 10}]
}
