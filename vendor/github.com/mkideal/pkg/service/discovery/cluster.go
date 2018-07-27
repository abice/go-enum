package discovery

import (
	"errors"
	"sync"

	"github.com/mkideal/pkg/math/random"
	"github.com/mkideal/pkg/netutil"
	"google.golang.org/grpc"
)

var (
	errNoAddr = errors.New("no available addr")
)

type GetResponse struct {
	Address grpc.Address
	Error   error

	request *GetRequest `json:"-"`
	waitCh  chan error  `json:"-"`
}

type GetRequest struct {
}

type SessionCreator func(addr string) netutil.Session

type clusterOptions struct {
	sessionCreator SessionCreator
}

type ClusterOption func(opts *clusterOptions)

func WithSessionCreator(sessionCreator SessionCreator) ClusterOption {
	return func(opts *clusterOptions) {
		opts.sessionCreator = sessionCreator
	}
}

type Cluster struct {
	opts        clusterOptions
	serviceName string
	balancer    grpc.Balancer
	addrs       []grpc.Address
	getCh       chan *GetResponse
	broadcastCh chan netutil.Packet
	quitCh      chan struct{}

	sessions map[string]netutil.Session
}

func NewCluster(serviceName string, balancer grpc.Balancer, opts ...ClusterOption) *Cluster {
	cluster := &Cluster{
		serviceName: serviceName,
		balancer:    balancer,
		getCh:       make(chan *GetResponse, 1024),
		broadcastCh: make(chan netutil.Packet, 1024),
		quitCh:      make(chan struct{}),
		sessions:    make(map[string]netutil.Session),
	}
	for _, opt := range opts {
		opt(&cluster.opts)
	}
	return cluster
}

func (cluster *Cluster) Startup() error {
	if err := cluster.balancer.Start(cluster.serviceName, grpc.BalancerConfig{}); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case addrs := <-cluster.balancer.Notify():
				cluster.addrs = addrs
				if cluster.opts.sessionCreator != nil {
					cluster.updateSessions()
				}
			case resp := <-cluster.getCh:
				if len(cluster.addrs) == 0 {
					resp.waitCh <- errNoAddr
				} else {
					// TODO: random select addr. is this okay?
					resp.Address = cluster.addrs[random.Intn(len(cluster.addrs), nil)]
					resp.waitCh <- nil
				}
			case packet := <-cluster.broadcastCh:
				for _, session := range cluster.sessions {
					session.Send(packet)
				}
			case <-cluster.quitCh:
				return
			}
		}
	}()
	return nil
}

func (cluster *Cluster) Shutdown(wg *sync.WaitGroup) {
	cluster.quitCh <- struct{}{}
	wg.Done()
}

func (cluster *Cluster) Get(req *GetRequest) *GetResponse {
	resp := &GetResponse{
		request: req,
		waitCh:  make(chan error),
	}
	cluster.getCh <- resp
	err := <-resp.waitCh
	resp.Error = err
	return resp
}

func (cluster *Cluster) Broadcast(packet netutil.Packet) {
	cluster.broadcastCh <- packet
}

func (cluster *Cluster) updateSessions() {
	open := make(map[string]bool)
	for _, addr := range cluster.addrs {
		open[addr.Addr] = true
		if oldSession, exist := cluster.sessions[addr.Addr]; !exist || oldSession.Closed() {
			cluster.newSession(addr.Addr)
		}
	}
	for k, s := range cluster.sessions {
		if !open[k] {
			delete(cluster.sessions, k)
			if s != nil {
				s.Quit()
			}
		}
	}
}

func (cluster *Cluster) newSession(addr string) {
	session := cluster.opts.sessionCreator(addr)
	if session == nil {
		return
	}
	cluster.sessions[addr] = session
	go session.Run(nil, nil)
}
