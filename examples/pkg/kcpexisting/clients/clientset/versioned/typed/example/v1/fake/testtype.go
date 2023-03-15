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

package fake

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kcp-dev/logicalcluster/v3"

	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"

	examplev1 "acme.corp/pkg/apis/example/v1"
	applyconfigurationsexamplev1 "acme.corp/pkg/generated/applyconfigurations/example/v1"
	examplev1client "acme.corp/pkg/generated/clientset/versioned/typed/example/v1"
	kcpexamplev1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v1"
)

var testTypesResource = schema.GroupVersionResource{Group: "example", Version: "v1", Resource: "testtypes"}
var testTypesKind = schema.GroupVersionKind{Group: "example", Version: "v1", Kind: "TestType"}

type testTypesClusterClient struct {
	*kcptesting.Fake
}

// Cluster scopes the client down to a particular cluster.
func (c *testTypesClusterClient) Cluster(clusterPath logicalcluster.Path) kcpexamplev1.TestTypesNamespacer {
	if clusterPath == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return &testTypesNamespacer{Fake: c.Fake, ClusterPath: clusterPath}
}

// List takes label and field selectors, and returns the list of TestTypes that match those selectors across all clusters.
func (c *testTypesClusterClient) List(ctx context.Context, opts metav1.ListOptions) (*examplev1.TestTypeList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewListAction(testTypesResource, testTypesKind, logicalcluster.Wildcard, metav1.NamespaceAll, opts), &examplev1.TestTypeList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &examplev1.TestTypeList{ListMeta: obj.(*examplev1.TestTypeList).ListMeta}
	for _, item := range obj.(*examplev1.TestTypeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested TestTypes across all clusters.
func (c *testTypesClusterClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewWatchAction(testTypesResource, logicalcluster.Wildcard, metav1.NamespaceAll, opts))
}

type testTypesNamespacer struct {
	*kcptesting.Fake
	ClusterPath logicalcluster.Path
}

func (n *testTypesNamespacer) Namespace(namespace string) examplev1client.TestTypeInterface {
	return &testTypesClient{Fake: n.Fake, ClusterPath: n.ClusterPath, Namespace: namespace}
}

type testTypesClient struct {
	*kcptesting.Fake
	ClusterPath logicalcluster.Path
	Namespace   string
}

func (c *testTypesClient) Create(ctx context.Context, testType *examplev1.TestType, opts metav1.CreateOptions) (*examplev1.TestType, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewCreateAction(testTypesResource, c.ClusterPath, c.Namespace, testType), &examplev1.TestType{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.TestType), err
}

func (c *testTypesClient) Update(ctx context.Context, testType *examplev1.TestType, opts metav1.UpdateOptions) (*examplev1.TestType, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewUpdateAction(testTypesResource, c.ClusterPath, c.Namespace, testType), &examplev1.TestType{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.TestType), err
}

func (c *testTypesClient) UpdateStatus(ctx context.Context, testType *examplev1.TestType, opts metav1.UpdateOptions) (*examplev1.TestType, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewUpdateSubresourceAction(testTypesResource, c.ClusterPath, "status", c.Namespace, testType), &examplev1.TestType{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.TestType), err
}

func (c *testTypesClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.Invokes(kcptesting.NewDeleteActionWithOptions(testTypesResource, c.ClusterPath, c.Namespace, name, opts), &examplev1.TestType{})
	return err
}

func (c *testTypesClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := kcptesting.NewDeleteCollectionAction(testTypesResource, c.ClusterPath, c.Namespace, listOpts)

	_, err := c.Fake.Invokes(action, &examplev1.TestTypeList{})
	return err
}

func (c *testTypesClient) Get(ctx context.Context, name string, options metav1.GetOptions) (*examplev1.TestType, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewGetAction(testTypesResource, c.ClusterPath, c.Namespace, name), &examplev1.TestType{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.TestType), err
}

// List takes label and field selectors, and returns the list of TestTypes that match those selectors.
func (c *testTypesClient) List(ctx context.Context, opts metav1.ListOptions) (*examplev1.TestTypeList, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewListAction(testTypesResource, testTypesKind, c.ClusterPath, c.Namespace, opts), &examplev1.TestTypeList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &examplev1.TestTypeList{ListMeta: obj.(*examplev1.TestTypeList).ListMeta}
	for _, item := range obj.(*examplev1.TestTypeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

func (c *testTypesClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.InvokesWatch(kcptesting.NewWatchAction(testTypesResource, c.ClusterPath, c.Namespace, opts))
}

func (c *testTypesClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*examplev1.TestType, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewPatchSubresourceAction(testTypesResource, c.ClusterPath, c.Namespace, name, pt, data, subresources...), &examplev1.TestType{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.TestType), err
}

func (c *testTypesClient) Apply(ctx context.Context, applyConfiguration *applyconfigurationsexamplev1.TestTypeApplyConfiguration, opts metav1.ApplyOptions) (*examplev1.TestType, error) {
	if applyConfiguration == nil {
		return nil, fmt.Errorf("applyConfiguration provided to Apply must not be nil")
	}
	data, err := json.Marshal(applyConfiguration)
	if err != nil {
		return nil, err
	}
	name := applyConfiguration.Name
	if name == nil {
		return nil, fmt.Errorf("applyConfiguration.Name must be provided to Apply")
	}
	obj, err := c.Fake.Invokes(kcptesting.NewPatchSubresourceAction(testTypesResource, c.ClusterPath, c.Namespace, *name, types.ApplyPatchType, data), &examplev1.TestType{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.TestType), err
}

func (c *testTypesClient) ApplyStatus(ctx context.Context, applyConfiguration *applyconfigurationsexamplev1.TestTypeApplyConfiguration, opts metav1.ApplyOptions) (*examplev1.TestType, error) {
	if applyConfiguration == nil {
		return nil, fmt.Errorf("applyConfiguration provided to Apply must not be nil")
	}
	data, err := json.Marshal(applyConfiguration)
	if err != nil {
		return nil, err
	}
	name := applyConfiguration.Name
	if name == nil {
		return nil, fmt.Errorf("applyConfiguration.Name must be provided to Apply")
	}
	obj, err := c.Fake.Invokes(kcptesting.NewPatchSubresourceAction(testTypesResource, c.ClusterPath, c.Namespace, *name, types.ApplyPatchType, data, "status"), &examplev1.TestType{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.TestType), err
}

func (c *testTypesClient) CreateField(ctx context.Context, testTypeName string, field *examplev1.Field, opts metav1.CreateOptions) (*examplev1.Field, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewCreateSubresourceAction(testTypesResource, c.ClusterPath, testTypeName, "field", c.Namespace, field), &examplev1.Field{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.Field), err
}

func (c *testTypesClient) UpdateField(ctx context.Context, testTypeName string, field *examplev1.Field, opts metav1.UpdateOptions) (*examplev1.Field, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewUpdateSubresourceAction(testTypesResource, c.ClusterPath, "field", c.Namespace, field), &examplev1.Field{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.Field), err
}

func (c *testTypesClient) GetField(ctx context.Context, testTypeName string, options metav1.GetOptions) (*examplev1.Field, error) {
	obj, err := c.Fake.Invokes(kcptesting.NewGetSubresourceAction(testTypesResource, c.ClusterPath, "field", c.Namespace, testTypeName), &examplev1.Field{})
	if obj == nil {
		return nil, err
	}
	return obj.(*examplev1.Field), err
}
