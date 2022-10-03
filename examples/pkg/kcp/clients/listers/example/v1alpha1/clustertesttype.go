//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by kcp code-generator. DO NOT EDIT.

package v1alpha1

import (
	kcpcache "github.com/kcp-dev/apimachinery/pkg/cache"
	"github.com/kcp-dev/logicalcluster/v2"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	examplev1alpha1 "acme.corp/pkg/apis/example/v1alpha1"
)

// ClusterTestTypeClusterLister can list ClusterTestTypes across all workspaces, or scope down to a ClusterTestTypeLister for one workspace.
type ClusterTestTypeClusterLister interface {
	List(selector labels.Selector) (ret []*examplev1alpha1.ClusterTestType, err error)
	Cluster(cluster logicalcluster.Name) ClusterTestTypeLister
}

type clusterTestTypeClusterLister struct {
	indexer cache.Indexer
}

// NewClusterTestTypeClusterLister returns a new ClusterTestTypeClusterLister.
func NewClusterTestTypeClusterLister(indexer cache.Indexer) *clusterTestTypeClusterLister {
	return &clusterTestTypeClusterLister{indexer: indexer}
}

// List lists all ClusterTestTypes in the indexer across all workspaces.
func (s *clusterTestTypeClusterLister) List(selector labels.Selector) (ret []*examplev1alpha1.ClusterTestType, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*examplev1alpha1.ClusterTestType))
	})
	return ret, err
}

// Cluster scopes the lister to one workspace, allowing users to list and get ClusterTestTypes.
func (s *clusterTestTypeClusterLister) Cluster(cluster logicalcluster.Name) ClusterTestTypeLister {
	return &clusterTestTypeLister{indexer: s.indexer, cluster: cluster}
}

type ClusterTestTypeLister interface {
	List(selector labels.Selector) (ret []*examplev1alpha1.ClusterTestType, err error)
	Get(name string) (*examplev1alpha1.ClusterTestType, error)
}

// clusterTestTypeLister can list all ClusterTestTypes inside a workspace.
type clusterTestTypeLister struct {
	indexer cache.Indexer
	cluster logicalcluster.Name
}

// List lists all ClusterTestTypes in the indexer for a workspace.
func (s *clusterTestTypeLister) List(selector labels.Selector) (ret []*examplev1alpha1.ClusterTestType, err error) {
	selectAll := selector == nil || selector.Empty()

	list, err := s.indexer.ByIndex(kcpcache.ClusterIndexName, kcpcache.ClusterIndexKey(s.cluster))
	if err != nil {
		return nil, err
	}

	for i := range list {
		obj := list[i].(*examplev1alpha1.ClusterTestType)
		if selectAll {
			ret = append(ret, obj)
		} else {
			if selector.Matches(labels.Set(obj.GetLabels())) {
				ret = append(ret, obj)
			}
		}
	}

	return ret, err
}

// Get retrieves the ClusterTestType from the indexer for a given workspace and name.
func (s *clusterTestTypeLister) Get(name string) (*examplev1alpha1.ClusterTestType, error) {
	key := kcpcache.ToClusterAwareKey(s.cluster.String(), "", name)
	obj, exists, err := s.indexer.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(examplev1alpha1.Resource("ClusterTestType"), name)
	}
	return obj.(*examplev1alpha1.ClusterTestType), nil
}
