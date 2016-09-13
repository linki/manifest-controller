package cluster

import (
	"os/exec"
	"strings"

	"github.com/linki/barn/manifest"

	log "github.com/Sirupsen/logrus"
)

type Cluster struct {
	url string
}

func New(url string) (*Cluster, error) {
	return &Cluster{url: url}, nil
}

func (c *Cluster) Apply(manifest manifest.Source) error {
	log.Info("Applying manifest")

	log.Debugf("\n%s", manifest.Manifest())

	err := c.applyCommand(manifest).Run()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cluster) applyCommand(manifest manifest.Source) *exec.Cmd {
	cmd := exec.Command("kubectl", "--server", c.url, "apply", "-f", "-")

	cmd.Stdin = strings.NewReader(manifest.Manifest())

	cmd.Stdout = log.StandardLogger().WriterLevel(log.DebugLevel)
	cmd.Stderr = log.StandardLogger().WriterLevel(log.ErrorLevel)

	return cmd
}
