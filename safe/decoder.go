package safe

import (
	"bytes"
	"fmt"
	"github.com/francoispqt/gojay"
	"strconv"
)

func DecodeInt(dec *gojay.Decoder) (int, error) {
	var raw = gojay.EmbeddedJSON{}
	if err := dec.EmbeddedJSON(&raw); err != nil {
		return 0, fmt.Errorf("failed to decoded embeded json %w",err)
	}
	if count := bytes.Count(raw, []byte(`"`)); count == 2 {
		raw = raw[1 : len(raw)-2]
	}
	if len(raw) == 0 {
		return 0, nil
	}
	return strconv.Atoi(string(raw))
}
