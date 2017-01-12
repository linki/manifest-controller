package manifest

import (
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Manifest struct {
	Content string

	client     *http.Client
	url        string
	lastUpdate time.Time

	sync.Mutex
}

func New(url string, client *http.Client) *Manifest {
	return &Manifest{url: url, client: client}
}

func (m *Manifest) Fetch() (bool, error) {
	log.Info("Downloading manifest")

	resp, err := m.client.Get(m.url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	manifest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	log.Debugf("\n%s", manifest)

	changed := m.Content != string(manifest)
	if !changed {
		log.Infof("Manifest unchanged")
	}

	m.Lock()
	{
		m.Content = string(manifest)
		m.lastUpdate = time.Now()
	}
	m.Unlock()

	return changed, nil
}

func (m *Manifest) Manifest() string {
	return m.Content
}
