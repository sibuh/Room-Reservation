package forms

import (
	"fmt"
	"net/url"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}
func (f *Form) Has(field string, fo url.Values) {
	x := fo.Get("field")
	fmt.Println(x)
	if x == "" {
		f.Errors.Add(field, "This field is mandatory")
	}
}
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
