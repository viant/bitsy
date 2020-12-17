package dec

import "github.com/francoispqt/gojay"

type Floats struct {
	Callback func(value float64)
}
// implement UnmarshalerJSONArray
func (s *Floats) UnmarshalJSONArray(dec *gojay.Decoder) error {

	value :=0.0
	if err := dec.Float(&value); err != nil {
		return err
	}
	s.Callback(value)
	return nil
}
