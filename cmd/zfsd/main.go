package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/google/go-microservice-helpers/server"
	"github.com/google/go-microservice-helpers/tracing"
	"github.com/steigr/zfsd/pkg/proto/zfs"
	"github.com/steigr/zfsd/pkg/proto/zpool"
	zfs_server "github.com/steigr/zfsd/pkg/zfs"
	zpool_server "github.com/steigr/zfsd/pkg/zpool"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	err := tracing.InitTracer(*serverhelpers.ListenAddress, "zfsd")
	if err != nil {
		glog.Fatalf("failed to init tracing interface: %v", err)
	}

	zpoolSrv := zpool_server.New()
	zfsSrv := zfs_server.New()

	grpcServer, _, err := serverhelpers.NewServer()
	if err != nil {
		glog.Fatalf("failed to init GRPC server: %v", err)
	}

	zpool.RegisterZPoolServer(grpcServer, &zpoolSrv)
	zfs.RegisterZfsServer(grpcServer, &zfsSrv)

	err = serverhelpers.ListenAndServe(grpcServer, nil)
	if err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}
