package server

import (
	"context"
	pb "github.com/steigr/zfsd/pkg/proto/zfs"
)

type Server struct{}

func New() Server {
	return Server{}
}

func (s Server) ListVolumes(ctx context.Context, in *pb.ListVolumesRequest) (out *pb.ListVolumesReply, err error) {
	return nil, nil
}