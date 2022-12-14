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

package v2

import (
	"context"

	kcpclient "github.com/kcp-dev/apimachinery/v2/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	examplev2 "acme.corp/pkg/apis/example/v2"
	examplev2client "acme.corp/pkg/generated/clientset/versioned/typed/example/v2"
)

// ClusterTestTypesClusterGetter has a method to return a ClusterTestTypeClusterInterface.
// A group's cluster client should implement this interface.
type ClusterTestTypesClusterGetter interface {
	ClusterTestTypes() ClusterTestTypeClusterInterface
}

// ClusterTestTypeClusterInterface can operate on ClusterTestTypes across all clusters,
// or scope down to one cluster and return a examplev2client.ClusterTestTypeInterface.
type ClusterTestTypeClusterInterface interface {
	Cluster(logicalcluster.Path) examplev2client.ClusterTestTypeInterface
	List(ctx context.Context, opts metav1.ListOptions) (*examplev2.ClusterTestTypeList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type clusterTestTypesClusterInterface struct {
	clientCache kcpclient.Cache[*examplev2client.ExampleV2Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *clusterTestTypesClusterInterface) Cluster(clusterPath logicalcluster.Path) examplev2client.ClusterTestTypeInterface {
	if clusterPath == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return c.clientCache.ClusterOrDie(clusterPath).ClusterTestTypes()
}

// List returns the entire collection of all ClusterTestTypes across all clusters.
func (c *clusterTestTypesClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*examplev2.ClusterTestTypeList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).ClusterTestTypes().List(ctx, opts)
}

// Watch begins to watch all ClusterTestTypes across all clusters.
func (c *clusterTestTypesClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).ClusterTestTypes().Watch(ctx, opts)
}
