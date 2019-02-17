package csv

/*
UnmarshalFunc is the method implemented by an object that can unmarshal a textual representation of the v.

UnmarshalFunc must be able to decode the form generated by MarshalFunc. UnmarshalFunc must copy the text if it wishes to retain the text after returning.
*/
type UnmarshalFunc func(v interface{}, text []byte) error