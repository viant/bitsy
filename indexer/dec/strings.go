package dec

import "github.com/francoispqt/gojay"

type Strings struct {
	Callback func(value string)
}
// implement UnmarshalerJSONArray
func (s *Strings) UnmarshalJSONArray(dec *gojay.Decoder) error {
	value := ""
	if err := dec.String(&value); err != nil {
		return err
	}
	s.Callback(value)
	return nil
}
