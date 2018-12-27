// +build linux,cgo

package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/steigr/zfsd/pkg/proto"
	"github.com/steigr/zfsd/pkg/server"
	"github.com/steigr/zfsd/pkg/util"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	grpcServer, err := util.NewServer()
	if err != nil {
		glog.Fatalf("failed to init GRPC server: %v", err)
	}

	zfs.RegisterDatasetServer(grpcServer, server.DatasetServer())
	zfs.RegisterPoolServer(grpcServer, server.PoolServer())

	err = util.ListenAndServe(grpcServer)
	if err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}
