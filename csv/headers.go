package csv

import "fmt"

type headerMap map[string]int

type headerList []string

func (headers headerList) ToMap() headerMap {
	headerMap := headerMap{}
	for i, header := range headers {
		headerMap[header] = i
	}
	return headerMap
}

func getHeaders(r csvReader, headers headerList) (headerMap, error) {

	if r == nil {
		return nil, fmt.Errorf("reader can't be nil")
	}

	if headers != nil {
		return headers.ToMap(), nil
	}

	headerBytes, err := r.Read()
	if err != nil {
		return nil, err
	}

	headermap := headerMap{}
	for i, headerByte := range headerBytes {
		headermap[string(headerByte)] = i
	}

	return headermap, nil
}
