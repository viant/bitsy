package dec

import "github.com/francoispqt/gojay"

type Floats struct {
	Items []float64
}
// implement UnmarshalerJSONArray
func (s *Floats) UnmarshalJSONArray(dec *gojay.Decoder) error {
	if len(s.Items) == 0 {
		s.Items = make([]float64, 0)
	}
	value :=0.0
	if err := dec.Float(&value); err != nil {
		return err
	}
	s.Items = append(s.Items, value)
	return nil
}
