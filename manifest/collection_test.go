package manifest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	otherManifestPath = "/some-other-path.yaml"
)

var (
	server2 *httptest.Server
)

func init() {
	responses := map[string]string{
		manifestPath:      "some manifest",
		otherManifestPath: "some other manifest",
	}

	server2 = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != manifestPath && r.URL.Path != otherManifestPath {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(responses[r.URL.Path]))
	}))
}

func TestCollection(t *testing.T) {
	var (
		manifest = New(server2.URL+manifestPath, http.DefaultClient)
		other    = New(server2.URL+otherManifestPath, http.DefaultClient)
	)

	collection := NewCollection(manifest, other)

	_, err := collection.Fetch()
	if err != nil {
		t.Errorf("%s", err)
	}

	actual := collection.Content

	expected := `some manifest

---

some other manifest`

	if actual != expected {
		t.Errorf("collection.Content == %q, want %q", actual, expected)
	}
}
