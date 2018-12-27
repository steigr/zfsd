package server

import (
	"context"
	libzfs "github.com/bicomsystems/go-libzfs"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
	"github.com/steigr/zfsd/pkg/proto"
)

type DatasetServerImpl struct {}

func DatasetServer() *DatasetServerImpl {
	return &DatasetServerImpl{}
}

func (s DatasetServerImpl) List(ctx context.Context, in *zfs.ListDatasetRequest) (out *zfs.ListDatasetReply, err error) {
	glog.V(9).Info("Dataset.List()")
	datasets, err := getAllDatasets(in.UserProperties)
	if err != nil {
		return nil, err
	}
	out = &zfs.ListDatasetReply{
		Datasets: datasets,
	}
	return out,err
}

func (s DatasetServerImpl) Get(ctx context.Context, in *zfs.GetDatasetRequest) (out *zfs.GetDatasetReply, err error) {
	glog.V(9).Info("Dataset.Get()")
	var libzfsDataset libzfs.Dataset
	libzfsDataset, err = libzfs.DatasetOpen(in.Path)
	if err != nil {
		return nil, err
	}
	defer libzfsDataset.Close()
	out = &zfs.GetDatasetReply{
		Dataset: toDatasetT(libzfsDataset,in.UserProperties),
	}
	return out, nil
}

func (s DatasetServerImpl) Create(ctx context.Context, in *zfs.CreateDatasetRequest) (out *zfs.CreateDatasetReply, err error) {
	glog.V(9).Info("Dataset.Create()")
	dataset, err := libzfs.DatasetCreate(in.Dataset.Path,libzfs.DatasetTypeFilesystem,make(map[libzfs.Prop]libzfs.Property,0))
	if err != nil {
		return nil, err
	}
	assignProperties(&dataset,in.Dataset.Properties)
	out = &zfs.CreateDatasetReply{
		Dataset: toDatasetT(dataset,make([]string,0)),
	}
	return out, nil
}

func (s DatasetServerImpl) Update(ctx context.Context, in *zfs.UpdateDatasetRequest) (out *zfs.UpdateDatasetReply, err error) {
	glog.V(9).Info("Dataset.Update()")
	spew.Dump(in)
	dataset, err := libzfs.DatasetOpen(in.Dataset.Path)
	if err != nil {
		return nil, err
	}
	defer dataset.Close()
	assignProperties(&dataset,in.Dataset.Properties)
	out = &zfs.UpdateDatasetReply{
		Dataset: toDatasetT(dataset,make([]string,0)),
	}
	return out, nil
}

func (s DatasetServerImpl) Delete(ctx context.Context, in *zfs.DeleteDatasetRequest) (out *zfs.DeleteDatasetReply, err error) {
	out = &zfs.DeleteDatasetReply{
		Success: false,
	}
	dataset, err := libzfs.DatasetOpen(in.Path)
	if err != nil {
		return out, err
	}
	defer dataset.Close()
	err = dataset.Destroy(true)
	if err != nil {
		return out, err
	}
	out.Success = true
	return out, nil
}

func getAllDatasets(userProperties []string) (datasets []*zfs.DatasetT, err error) {
	libzfsDatasets, err := libzfs.DatasetOpenAll()
	if err != nil {
		return nil, err
	}
	defer libzfs.DatasetCloseAll(libzfsDatasets)
	datasets = make([]*zfs.DatasetT,len(libzfsDatasets))
	for idx, libzfsDataset := range libzfsDatasets {
		datasets[idx] = toDatasetT(libzfsDataset, userProperties)
	}
	return datasets, nil
}

func toDatasetT(in libzfs.Dataset,userProperties []string) (out *zfs.DatasetT) {
	out = &zfs.DatasetT{}
	out.Path = in.Properties[libzfs.DatasetPropName].Value
	for idx, prop := range in.Properties {
		out.Properties = append(out.Properties,&zfs.Property{
			Key:                  libzfs.DatasetPropertyToName(libzfs.Prop(idx)),
			Source:               prop.Source,
			Value:                prop.Value,
		})
	}
	if len(userProperties) > 0 {
		for _, name := range userProperties {
			prop, err := in.GetUserProperty(name)
			if err != nil {
				continue
			}
			out.Properties = append(out.Properties,&zfs.Property{
				Key: name,
				Source:  prop.Source,
				Value: prop.Value,

			})
		}
	}
	if len(in.Children) > 0 {
		for _, child := range in.Children {
			out.Datasets = append(out.Datasets,toDatasetT(child,userProperties))
		}
	}
	return out
}

func assignProperties(dataset *libzfs.Dataset, inProperties []*zfs.Property) (err error) {
	path, err := dataset.Path()
	if err != nil {
		return err
	}
	for _, inProperty := range inProperties {
		propIdx, isset :=  datasetProperties[inProperty.Key]
		if isset {
			propType := libzfs.Prop(propIdx)
			glog.V(9).Info("assignProperties() ",path," property ",libzfs.DatasetPropertyToName(propType)," = ",inProperty.Value, " assigned")
			dataset.SetProperty(propType,inProperty.Value)
		} else {
			glog.V(9).Info("assignProperties() ",path," user-property ",inProperty.Key," = ",inProperty.Value, " assigned")
			dataset.SetUserProperty(inProperty.Key,inProperty.Value)
		}
	}
	return nil
}

func init() {
	collectDatasetProperties()
}

var datasetProperties map[string]int

func collectDatasetProperties() {
	datasetProperties = make(map[string]int,int(libzfs.DatasetNumProps))
	for propIdx := 0; propIdx < int(libzfs.DatasetNumProps); propIdx++ {
		datasetProperties[libzfs.DatasetPropertyToName(libzfs.Prop(propIdx))] = propIdx
	}
}
