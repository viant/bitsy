package indexer

import (
	"strings"
	"time"
)

const (
	timePathVar    = "$TimePath"
	pathTimeLayout = "2006/01/02/03"
)

func ExpandURL(URL string, time time.Time) string {
	if count := strings.Count(URL, timePathVar); count > 0 {
		URL = strings.Replace(URL, timePathVar, time.Format(pathTimeLayout), count)
	}
	return URL
}
