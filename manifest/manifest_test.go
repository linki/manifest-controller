package manifest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

const (
	manifestPath = "/some-path.yaml"
)

var (
	server  *httptest.Server
	content string
)

func init() {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != manifestPath {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(content))
	}))
}

func TestFetch(t *testing.T) {
	manifest := New(server.URL+manifestPath, http.DefaultClient)

	content = "some manifest"

	_, err := manifest.Fetch()
	if err != nil {
		t.Errorf("%s", err)
	}

	if manifest.Content != content {
		t.Errorf("manifest.Content == %q, want %q", manifest.Content, content)
	}
}

func TestChanged(t *testing.T) {
	manifest := New(server.URL+manifestPath, http.DefaultClient)

	content = "some manifest"

	changed, err := manifest.Fetch()
	if err != nil {
		t.Errorf("%s", err)
	}

	if changed != true {
		t.Errorf("changed == %t, want %t", changed, true)
	}

	changed, err = manifest.Fetch()
	if err != nil {
		t.Errorf("%s", err)
	}

	if changed != false {
		t.Errorf("changed == %t, want %t", changed, true)
	}

	content = "some other manifest"

	changed, err = manifest.Fetch()
	if err != nil {
		t.Errorf("%s", err)
	}

	if changed != true {
		t.Errorf("changed == %t, want %t", changed, true)
	}
}

func TestAuthHeader(t *testing.T) {
	var authHeader string

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
	}))

	client := oauth2.NewClient(context.TODO(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "foo"}),
	)

	manifest := New(server.URL+manifestPath, client)

	_, err := manifest.Fetch()
	if err != nil {
		t.Errorf("%s", err)
	}

	if authHeader != "Bearer foo" {
		t.Errorf("authHeader == %q, want %q", authHeader, "Bearer foo")
	}
}
