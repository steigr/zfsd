package util

import "github.com/steigr/zfsd/pkg/proto"

func NewPoolClient(endpoint string) (pools zfs.PoolClient, err error) {
	if err = Connect(endpoint); err != nil {
		return nil, err
	}
	pools = zfs.NewPoolClient(Client)
	return pools, nil
}
