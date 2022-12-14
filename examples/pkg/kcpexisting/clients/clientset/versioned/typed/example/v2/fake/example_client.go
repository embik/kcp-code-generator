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
	"github.com/kcp-dev/logicalcluster/v3"

	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"
	"k8s.io/client-go/rest"

	examplev2 "acme.corp/pkg/generated/clientset/versioned/typed/example/v2"
	kcpexamplev2 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v2"
)

var _ kcpexamplev2.ExampleV2ClusterInterface = (*ExampleV2ClusterClient)(nil)

type ExampleV2ClusterClient struct {
	*kcptesting.Fake
}

func (c *ExampleV2ClusterClient) Cluster(clusterPath logicalcluster.Path) examplev2.ExampleV2Interface {
	if clusterPath == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
	return &ExampleV2Client{Fake: c.Fake, ClusterPath: clusterPath}
}

func (c *ExampleV2ClusterClient) TestTypes() kcpexamplev2.TestTypeClusterInterface {
	return &testTypesClusterClient{Fake: c.Fake}
}

func (c *ExampleV2ClusterClient) ClusterTestTypes() kcpexamplev2.ClusterTestTypeClusterInterface {
	return &clusterTestTypesClusterClient{Fake: c.Fake}
}

var _ examplev2.ExampleV2Interface = (*ExampleV2Client)(nil)

type ExampleV2Client struct {
	*kcptesting.Fake
	ClusterPath logicalcluster.Path
}

func (c *ExampleV2Client) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}

func (c *ExampleV2Client) TestTypes(namespace string) examplev2.TestTypeInterface {
	return &testTypesClient{Fake: c.Fake, ClusterPath: c.ClusterPath, Namespace: namespace}
}

func (c *ExampleV2Client) ClusterTestTypes() examplev2.ClusterTestTypeInterface {
	return &clusterTestTypesClient{Fake: c.Fake, ClusterPath: c.ClusterPath}
}
