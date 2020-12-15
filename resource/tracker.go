package resource

import (
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"golang.org/x/net/context"
	"sync"
	"time"
)

type Tracker struct {
	baseURL        string
	assets         Assets
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

func (m *Tracker) hasChanges(assets []storage.Object) bool {
	if len(assets) != len(m.assets) {
		return true
	}
	for _, asset := range assets {
		if asset.IsDir() {
			continue
		}
		mAsset, ok := m.assets[asset.URL()]
		if !ok {
			return true
		}
		if !mAsset.ModTime().Equal(asset.ModTime()) {
			return true
		}
	}
	return false

}

//Notify returns true if resource under base URL have changed
func (m *Tracker) Notify(ctx context.Context, fs afs.Service, callback func(URL string, operation Operation)) error {
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
	assets := NewAssets(resources)

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if len(m.assets) == 0 {
		m.assets = make(map[string]storage.Object)
	}

	m.assets.Added(assets, func(object storage.Object) {
		callback(object.URL(), OperationAdded)
	})
	m.assets.Modified(assets, func(object storage.Object) {
		callback(object.URL(), OperationModified)
	})
	m.assets.Deleted(assets, func(object storage.Object) {
		callback(object.URL(), OperationDeleted)
	})
	return nil
}

func New(baseURL string, checkFrequency time.Duration) *Tracker {
	if checkFrequency == 0 {
		checkFrequency = time.Minute
	}
	return &Tracker{
		checkFrequency: checkFrequency,
		mutex:          sync.Mutex{},
		baseURL:        baseURL,
		assets:         make(map[string]storage.Object),
	}
}
