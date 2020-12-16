package dec

import "github.com/francoispqt/gojay"

type Strings struct {
	Items []string
}
// implement UnmarshalerJSONArray
func (s *Strings) UnmarshalJSONArray(dec *gojay.Decoder) error {
	if len(s.Items) == 0 {
		s.Items = make([]string, 0)
	}
	value := ""
	if err := dec.String(&value); err != nil {
		return err
	}
	s.Items = append(s.Items, value)
	return nil
}
