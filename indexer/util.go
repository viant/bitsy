package indexer

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	timePathVar    = "$TimePath"
	pathTimeLayout = "2006/01/02/03"
	uuidVar        = "$UUID"
)

func ExpandURL(URL string, time time.Time) string {
	if count := strings.Count(URL, uuidVar); count > 0 {
		URL = strings.Replace(URL, uuidVar, uuid.New().String(), count)
	}
	if count := strings.Count(URL, timePathVar); count > 0 {
		URL = strings.Replace(URL, timePathVar, time.Format(pathTimeLayout), count)
	}
	return URL
}
