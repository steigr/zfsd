package util

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	// ListenAddress is the grpc listen address
	ListenAddress = flag.String("listen", "tcp://0.0.0.0:1737", "GRPC listen address")
)

// ListenAndServe starts grpc server
func ListenAndServe(grpcServer *grpc.Server) error {
	listenAddressTuple := strings.Split(*ListenAddress,"://")
	lis, err := net.Listen(listenAddressTuple[0], listenAddressTuple[1])
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	glog.Warningf("serving INSECURE on %v", *ListenAddress)
	err = grpcServer.Serve(lis)
	return fmt.Errorf("failed to serve: %v", err)
}

// NewServer creates a new GRPC server stub
func NewServer() (*grpc.Server, error) {
	var grpcServer *grpc.Server

	grpcServer = grpc.NewServer()

	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)

	return grpcServer, nil
}