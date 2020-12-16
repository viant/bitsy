package dec

import (
	"github.com/francoispqt/gojay"
	"github.com/viant/bitsy/safe"
)

type Ints struct {
	Items              []int
	AllowQuotedNumbers bool
}

// implement UnmarshalerJSONArray
func (s *Ints) UnmarshalJSONArray(dec *gojay.Decoder) (err error) {
	if len(s.Items) == 0 {
		s.Items = make([]int, 0)
	}
	value := 0
	if s.AllowQuotedNumbers {
		if value, err = safe.DecodeInt(dec); err != nil {
			return err
		}
	}
	if err := dec.Int(&value); err != nil {
		return err
	}
	s.Items = append(s.Items, value)
	return nil
}
