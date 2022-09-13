package sd

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
)

// Service discovery registry
type Registry interface {
	// Returns list of nodes for specified service
	Get(ctx context.Context, service string, tags []string, waitIdx uint64) ([]string, uint64, error)
}

type registry struct {
	consul     *api.Client
	datacenter string
	cluster    string
}

// Creates service discovery registry
func NewRegistry(consul *api.Client, datacenter string, cluster string) Registry {
	return &registry{
		consul:     consul,
		datacenter: datacenter,
		cluster:    cluster,
	}
}

// Returns list of nodes for specified service
func (r *registry) Get(ctx context.Context, service string, tags []string, waitIdx uint64) ([]string, uint64, error) {
	opts := &api.QueryOptions{
		Datacenter: r.datacenter,
		WaitIndex:  waitIdx,
	}
	opts = opts.WithContext(ctx)

	entries, meta, err := r.consul.Health().Service(service, r.cluster, true, opts)
	if err != nil {
		return nil, 0, err
	}

	if len(tags) > 0 {
		entries = r.filterEntries(entries, tags)
	}

	return r.parseEntries(entries), meta.LastIndex, nil
}

// Filter out entries that doesn't match to specified tags
func (r *registry) filterEntries(entries []*api.ServiceEntry, tags []string) (res []*api.ServiceEntry) {
EntriesLoop:
	for _, entry := range entries {
		for _, required := range tags {
			var found bool
			for _, tag := range entry.Service.Tags {
				if tag == required {
					found = true
					break
				}
			}

			if !found {
				continue EntriesLoop
			}
		}
		res = append(res, entry)
	}
	return
}

// Parse list of entries and returns list nodes
func (r *registry) parseEntries(entries []*api.ServiceEntry) []string {
	nodes := make([]string, len(entries))
	for i, entry := range entries {
		host := entry.Node.Address
		if entry.Service.Address != "" && entry.Service.Address != "127.0.0.1" && entry.Service.Address != "localhost" {
			host = entry.Service.Address
		}

		nodes[i] = fmt.Sprintf("%s:%d", host, entry.Service.Port)
	}

	return nodes
}
