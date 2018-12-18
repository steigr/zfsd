package zpool

import (
	"context"
	libzfs "github.com/bicomsystems/go-libzfs"
	pb "github.com/steigr/zfsd/pkg/proto/zpool"
)

type Server struct{}

func New() Server {
	return Server{}
}

func (s Server) ListPools(ctx context.Context, in *pb.ListPoolsRequest) (out *pb.ListPoolsReply, err error) {
	var (
		pools   []libzfs.Pool
		pbPools []*pb.Pool
	)

	out = &pb.ListPoolsReply{}

	if pools, err = libzfs.PoolOpenAll(); err != nil {
		return nil, err
	}

	for _, pool := range pools {
		pbPools = append(pbPools, &pb.Pool{
			Properties: []*pb.Pool_Property{
				{
					Value:  func() string { n, _ := pool.Name(); return n }(),
					Source: "none",
				},
			},
		})
	}

	out.Pools = pbPools

	return out, nil
}