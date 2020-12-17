package dec

import "github.com/francoispqt/gojay"

type Bools struct {
	Callback func(value bool)
}
// implement UnmarshalerJSONArray
func (s *Bools) UnmarshalJSONArray(dec *gojay.Decoder) error {

	value :=false
	if err := dec.Bool(&value); err != nil {
		return err
	}
	s.Callback(value)
	return nil
}
