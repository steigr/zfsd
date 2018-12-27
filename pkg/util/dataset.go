package util

import "github.com/steigr/zfsd/pkg/proto"

func NewDatasetClient(endpoint string) (datasets zfs.DatasetClient, err error) {
	if err = Connect(endpoint); err != nil {
		return nil, err
	}
	datasets = zfs.NewDatasetClient(Client)
	return datasets, nil
}

