# go-csv
Go CSV that uses reflect like JSON

```go
package main

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

var blob = []byte(`"Name","Age"
"Bob","12"
"Sally","13"
"Alice","10"
`)

func main() {
	var records []Record
	
	if err := csv.Unmarshal(&records, nil, blob); err != nil {
		log.Fatal(err)
	}

	fmt.Println(records)
	// Output:
	// [{Bob 12} {Sally 13} {Alice 10}]
}
```
Try on go playground https://play.golang.org/p/OS6P1e_Gs3s
