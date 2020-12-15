package resource

import (
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"golang.org/x/net/context"
	"path"
	"sync"
	"time"
)

type Tracker struct {
	baseURL        string
	resources      map[string]time.Time // map between resource url and last modification time
	mutex          sync.Mutex
	checkFrequency time.Duration
	nextCheck      time.Time
}

func (m *Tracker) isCheckDue(now time.Time) bool {
	if m.nextCheck.IsZero() || now.After(m.nextCheck) {
		m.nextCheck = now.Add(m.checkFrequency)
		return true
	}
	return false
}

func (m *Tracker) hasChanges(routes []storage.Object) bool {
	if len(routes) != len(m.resources) {
		return true
	}
	for _, route := range routes {
		if route.IsDir() {
			continue
		}
		modTime, ok := m.resources[route.URL()]
		if !ok {
			return true
		}
		if !modTime.Equal(route.ModTime()) {
			return true
		}
	}
	return false

}

//Notify returns true if resource under base URL have changed
func (m *Tracker) Notify(ctx context.Context, fs afs.Service, callback func(URL string, operation int)) error {
	if m.baseURL == "" {
		return nil
	}
	if !m.isCheckDue(time.Now()) {
		return nil
	}

	resources, err := fs.List(ctx, m.baseURL, option.NewRecursive(true))
	if err != nil {
		return errors.Wrapf(err, "failed to load rules %v", m.baseURL)
	}
	if !m.hasChanges(resources) {
		return nil
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.resources = make(map[string]time.Time)
	for _, route := range resources {
		if route.IsDir() || !(path.Ext(route.Name()) == ".json" || path.Ext(route.Name()) == ".yaml") {
			continue
		}
		m.resources[route.URL()] = route.ModTime()
	}
	return nil
}

func New(baeURL string, checkFrequency time.Duration) *Tracker {
	if checkFrequency == 0 {
		checkFrequency = time.Minute
	}
	return &Tracker{
		checkFrequency: checkFrequency,
		mutex:          sync.Mutex{},
		baseURL:        baeURL,
		resources:      make(map[string]time.Time),
	}
}
