package dec

import (
	"github.com/francoispqt/gojay"
	"strconv"
)

type Ints struct {
	Items    []int
	IsQuoted bool
}

// implement UnmarshalerJSONArray
func (s *Ints) UnmarshalJSONArray(dec *gojay.Decoder) (err error) {
	if len(s.Items) == 0 {
		s.Items = make([]int, 0)
	}
	value := 0
	if s.IsQuoted {
		text := ""
		if err = dec.String(&text); err != nil {
			return err
		}
		if value, err = strconv.Atoi(text); err != nil {
			return err
		}
		s.Items = append(s.Items, value)
		return nil

	}
	if err := dec.Int(&value); err != nil {
		return err
	}
	s.Items = append(s.Items, value)
	return nil
}
