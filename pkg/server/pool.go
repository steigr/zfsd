package server

import (
	"context"
	libzfs "github.com/bicomsystems/go-libzfs"
	"github.com/golang/glog"
	"github.com/steigr/zfsd/pkg/proto"
)

type PoolServerImpl struct {}

func PoolServer() *PoolServerImpl {
	return &PoolServerImpl{}
}

func (s PoolServerImpl) List(ctx context.Context, in *zfs.ListPoolRequest) (out *zfs.ListPoolReply, err error) {
	glog.V(9).Info("Pool.List()")
	out = &zfs.ListPoolReply{}
	pools, err := getAllPools()
	if err != nil {
		return nil, err
	} else {
		out.Pools = pools
	}
	return out,err
}

func (s PoolServerImpl) Get(ctx context.Context, in *zfs.GetPoolRequest) (out *zfs.GetPoolReply, err error) {
	glog.V(9).Info("Pool.Get()")
	out = &zfs.GetPoolReply{}
	pool, err := getPool(in.Name)
	if err != nil {
		return nil, err
	} else {
		out.Pool = pool
	}
	return out,err
}

func (s PoolServerImpl) Create(ctx context.Context, in *zfs.CreatePoolRequest) (out *zfs.CreatePoolReply, err error) {
	glog.V(9).Info("Pool.Create()")
	out = &zfs.CreatePoolReply{}
	return out,err
}

func (s PoolServerImpl) Delete(ctx context.Context, in *zfs.DeletePoolRequest) (out *zfs.DeletePoolReply, err error) {
	glog.V(9).Info("Pool.Delete()")
	out = &zfs.DeletePoolReply{}
	return out,err
}

func (s PoolServerImpl) Update(ctx context.Context, in *zfs.UpdatePoolRequest) (out *zfs.UpdatePoolReply, err error) {
	glog.V(9).Info("Pool.Update()")
	out = &zfs.UpdatePoolReply{}
	return out,err
}

func getAllPools() (pools []*zfs.PoolT, err error) {
	libzfsPools, err := libzfs.PoolOpenAll()
	if err != nil {
		return pools, err
	}
	pools = make([]*zfs.PoolT,len(libzfsPools))
	for idx, libzfsPool := range libzfsPools {
		pools[idx] = toPoolT(libzfsPool)
	}
	return pools, err
}

func getPool(poolName string) (pool *zfs.PoolT, err error) {
	if libzfsPool, err := libzfs.PoolOpen(poolName); err != nil {
		return pool, err
	} else {
		pool = toPoolT(libzfsPool)
	}
	return pool, nil
}

func toPoolT(in libzfs.Pool) (out *zfs.PoolT) {
	out = &zfs.PoolT{}
	out.Name = in.Properties[libzfs.PoolPropName].Value
	for idx, prop := range in.Properties {
		out.Properties = append(out.Properties,&zfs.Property{
			Key:                  libzfs.PoolPropertyToName(libzfs.Prop(idx)),
			Source:               prop.Source,
			Value:                prop.Value,
		})
	}
	return out
}

func init() {
	collectPoolProperties()
}

var poolProperties map[string]int

func collectPoolProperties() {
	poolProperties = make(map[string]int,int(libzfs.PoolNumProps))
	for propIdx := 0; propIdx < int(libzfs.PoolNumProps); propIdx++ {
		poolProperties[libzfs.PoolPropertyToName(libzfs.Prop(propIdx))] = propIdx
	}
}
