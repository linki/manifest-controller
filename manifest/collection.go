package manifest

import (
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type ManifestCollection struct {
	manifests  []Source
	Content    string
	lastUpdate time.Time

	sync.Mutex
}

func NewCollection(manifests ...Source) *ManifestCollection {
	return &ManifestCollection{
		manifests: manifests,
	}
}

func (mc *ManifestCollection) Fetch() (bool, error) {
	contents := make([]string, 0, len(mc.manifests))

	for _, m := range mc.manifests {
		_, err := m.Fetch()
		if err != nil {
			return false, err
		}

		contents = append(contents, m.Manifest())
	}

	manifest := strings.Join(contents, "\n\n---\n\n")

	log.Debugf("\n%s", manifest)

	changed := mc.Content != manifest
	if !changed {
		log.Infof("Manifest collection unchanged")
	}

	mc.Lock()
	{
		mc.Content = manifest
		mc.lastUpdate = time.Now()
	}
	mc.Unlock()

	return changed, nil
}

func (mc *ManifestCollection) Manifest() string {
	return mc.Content
}
