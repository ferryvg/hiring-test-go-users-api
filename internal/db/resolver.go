package db

import (
	"context"
	"sync"
	"time"

	"github.com/ferryvg/hiring-test-go-users-api/internal/core/sd"
	"github.com/sirupsen/logrus"
)

// Database nodes resolver.
type Resolver interface {

	// Initialize resolver.
	Init() error

	// Shutdowns resolver and frees used resources.
	Shutdown()
}

// Resolver config.
type ResolverConf struct {
	// Service name.
	Service string

	// Service tags.
	Tags []string
}

type resolverImpl struct {
	registry sd.Registry
	connList ConnList
	logger   logrus.FieldLogger
	conf     *ResolverConf
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	waitIdx  uint64
}

func NewResolver(registry sd.Registry, connList ConnList, logger logrus.FieldLogger, conf *ResolverConf) Resolver {
	ctx, cancel := context.WithCancel(context.Background())

	return &resolverImpl{
		registry: registry,
		connList: connList,
		logger:   logger,
		conf:     conf,
		ctx:      ctx,
		cancel:   cancel,
		wg:       new(sync.WaitGroup),
	}
}

func (r *resolverImpl) Init() error {
	if err := r.lookup(); err != nil {
		return err
	}

	r.wg.Add(1)

	go r.watch()

	return nil
}

func (r *resolverImpl) Shutdown() {
	r.cancel()
	r.wg.Wait()
}

func (r *resolverImpl) watch() {
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return

		default:
			err := r.lookup()
			if err != nil && r.ctx.Err() == nil {
				fields := logrus.Fields{
					"service": r.conf.Service,
					"tags":    r.conf.Tags,
				}

				r.logger.WithError(err).WithFields(fields).Error("Failed to lookup nodes")

				time.Sleep(time.Second)
			}
		}
	}
}

func (r *resolverImpl) lookup() error {
	nodes, waitIdx, err := r.registry.Get(r.ctx, r.conf.Service, r.conf.Tags, r.waitIdx)
	if err != nil {
		return err
	}

	r.waitIdx = waitIdx
	r.connList.SetNodes(nodes)

	return nil
}
