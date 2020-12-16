package dec

import "github.com/francoispqt/gojay"

type Bools struct {
	Items []bool

}
// implement UnmarshalerJSONArray
func (s *Bools) UnmarshalJSONArray(dec *gojay.Decoder) error {
	if len(s.Items) == 0 {
		s.Items = make([]bool, 0)
	}
	value :=false
	if err := dec.Bool(&value); err != nil {
		return err
	}
	s.Items = append(s.Items, value)
	return nil
}
