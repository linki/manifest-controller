package controller

import (
	"sync"
	"time"

	"github.com/linki/barn/cluster"
	"github.com/linki/barn/manifest"

	log "github.com/Sirupsen/logrus"
)

type controller struct {
	source   manifest.Source
	cluster  *cluster.Cluster
	updates  chan manifest.Source
	interval time.Duration

	sync.WaitGroup
}

type Options struct {
	Interval time.Duration
}

func New(source manifest.Source, cluster *cluster.Cluster, options *Options) *controller {
	return &controller{
		source:   source,
		cluster:  cluster,
		updates:  make(chan manifest.Source),
		interval: options.Interval,
	}
}

func (c *controller) Run() {
	c.monitorSource()
	c.applyUpdates()

	c.Wait()
}

func (c *controller) monitorSource() {
	c.Add(1)

	go func() {
		defer c.Done()

		for {
			changed, err := c.source.Fetch()
			if err != nil {
				log.Fatal(err)
			}

			if changed {
				c.updates <- c.source
			}

			log.Debugf("Sleeping for %s", c.interval)
			time.Sleep(c.interval)
		}
	}()
}

func (c *controller) applyUpdates() {
	c.Add(1)

	go func() {
		defer c.Done()

		for {
			err := c.cluster.Apply(<-c.updates)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}
