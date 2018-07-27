package discovery

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
)

// 服务地址
type Address struct {
	// 服务地址
	Addr string
	// 元数据
	Metadata string
}

// 服务发现
type Discovery struct {
	EtcdEndpoints string `json:",omitempty"` // etcd 端地址列表,以逗号分割

	// etcd 客户端
	client   *clientv3.Client         `json:"-"`
	resolver *etcdnaming.GRPCResolver `json:"-"`
}

func (dc *Discovery) Client() *clientv3.Client {
	return dc.client
}

func (dc *Discovery) SetCommandLineFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&dc.EtcdEndpoints, "etcd", os.Getenv("ETCD_ENDPOINTS"), "etcd endpoints, default=${ETCD_ENDPOINTS}")
}

func (dc *Discovery) Init() error {
	endpoints := strings.Split(dc.EtcdEndpoints, ",")
	var err error
	dc.client, err = clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		return err
	}
	dc.resolver = &etcdnaming.GRPCResolver{Client: dc.client}
	return nil
}

func (dc *Discovery) NewBalancer() grpc.Balancer {
	return grpc.RoundRobin(&etcdnaming.GRPCResolver{Client: dc.client})
}

func (dc *Discovery) Register(service string, address Address, opts ...clientv3.OpOption) error {
	return dc.resolver.Update(context.TODO(), service, naming.Update{
		Op:       naming.Add,
		Addr:     address.Addr,
		Metadata: address.Metadata,
	}, opts...)
}

func (dc *Discovery) Unregister(service string, addr string) error {
	return dc.resolver.Update(context.TODO(), service, naming.Update{
		Op:   naming.Delete,
		Addr: addr,
	})
}

func (dc *Discovery) Dial(service string) (*grpc.ClientConn, error) {
	b := dc.NewBalancer()
	conn, err := grpc.Dial(service, grpc.WithBalancer(b), grpc.WithInsecure())
	if err != nil {
		err = grpcDialError{service: service, err: err}
	}
	return conn, err
}

type grpcDialError struct {
	service string
	err     error
}

func (e grpcDialError) Error() string {
	return fmt.Sprintf("dial service %s: %v", e.service, e.err)
}
