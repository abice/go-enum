package discovery

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const DefaultTTL = time.Second * 30

type TTL struct {
	lease clientv3.Lease
	id    clientv3.LeaseID
}

func (ttl *TTL) ID() clientv3.LeaseID   { return ttl.id }
func (ttl *TTL) Opt() clientv3.OpOption { return clientv3.WithLease(ttl.id) }

func (ttl *TTL) Update() {
	ttl.lease.KeepAliveOnce(context.TODO(), ttl.id)
}

func Interval(discovery Discovery, f func(*TTL), duration time.Duration) (*TTL, error) {
	lease := clientv3.NewLease(discovery.client)
	resp, err := lease.Grant(context.TODO(), int64(duration/time.Second))
	if err != nil {
		return nil, err
	}
	ttl := &TTL{lease: lease, id: resp.ID}
	f(ttl)
	go func() {
		for range time.Tick(duration / 2) {
			f(ttl)
			ttl.Update()
		}
	}()
	return ttl, nil
}
