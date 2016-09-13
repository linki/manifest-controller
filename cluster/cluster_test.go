package cluster

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/linki/barn/manifest"
)

func TestApplyCommand(t *testing.T) {
	content := "some manifest"

	manifest := manifest.New("http://some-manifest", http.DefaultClient)
	manifest.Content = content

	cluster, err := New("http://some-url")
	if err != nil {
		t.Errorf("%s", err)
	}

	cmd := cluster.applyCommand(manifest)

	args := []string{
		"kubectl",
		"--server",
		"http://some-url",
		"apply",
		"-f",
		"-",
	}

	if len(cmd.Args) != len(args) {
		t.Errorf("len(cmd.Args) == %d, want %d", len(cmd.Args), len(args))
	}

	for i, _ := range args {
		if cmd.Args[i] != args[i] {
			t.Errorf("cmd.Args[%d] == %q, want %q", i, cmd.Args[i], args[i])
		}
	}

	stdin, err := ioutil.ReadAll(cmd.Stdin)
	if err != nil {
		t.Errorf("%s", err)
	}

	if !bytes.Equal(stdin, []byte(content)) {
		t.Errorf("cmd.Stdin == %s, want %s", stdin, content)
	}
}
