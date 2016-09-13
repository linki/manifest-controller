package main

import (
	"context"
	"time"

	"github.com/linki/barn/cluster"
	"github.com/linki/barn/controller"
	"github.com/linki/barn/manifest"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/oauth2"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	defaultCluster  = "http://127.0.0.1:8001"
	defaultInterval = 1 * time.Minute
)

var params struct {
	sources  []string
	cluster  string
	token    string
	interval time.Duration
	debug    bool
}

var version = "Unknown"

func init() {
	kingpin.Flag("source", "List of sources to watch (can be any OAuth2 protected HTTP resource)").StringsVar(&params.sources)
	kingpin.Flag("cluster", "The cluster to connect to (there's no means of authentication, use `kubectl proxy`)").Default(defaultCluster).StringVar(&params.cluster)
	kingpin.Flag("token", "An optional Bearer token sent with the request").StringVar(&params.token)
	kingpin.Flag("interval", "Interval in Duration format, e.g. 60s.").Short('i').Default(defaultInterval.String()).DurationVar(&params.interval)
	kingpin.Flag("debug", "Enable debug logging").BoolVar(&params.debug)
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	if params.debug {
		log.SetLevel(log.DebugLevel)
	}

	client := oauth2.NewClient(context.TODO(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: params.token}),
	)

	manifests := make([]*manifest.Manifest, 0)

	for _, s := range params.sources {
		manifests = append(manifests, manifest.New(s, client))
	}

	source := manifest.NewCollection(manifests...)

	cluster, err := cluster.New(params.cluster)
	if err != nil {
		log.Fatal(err)
	}

	ctrl := controller.New(source, cluster, &controller.Options{Interval: params.interval})
	ctrl.Run()
}
