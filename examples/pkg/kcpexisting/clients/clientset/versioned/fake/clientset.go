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
	"github.com/kcp-dev/logicalcluster/v2"

	kcpfakediscovery "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/discovery/fake"
	kcptesting "github.com/kcp-dev/client-go/third_party/k8s.io/client-go/testing"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"

	client "acme.corp/pkg/generated/clientset/versioned"
	clientscheme "acme.corp/pkg/generated/clientset/versioned/scheme"
	examplev1 "acme.corp/pkg/generated/clientset/versioned/typed/example/v1"
	examplev1alpha1 "acme.corp/pkg/generated/clientset/versioned/typed/example/v1alpha1"
	examplev1beta1 "acme.corp/pkg/generated/clientset/versioned/typed/example/v1beta1"
	examplev2 "acme.corp/pkg/generated/clientset/versioned/typed/example/v2"
	example3v1 "acme.corp/pkg/generated/clientset/versioned/typed/example3/v1"
	existinginterfacesv1 "acme.corp/pkg/generated/clientset/versioned/typed/existinginterfaces/v1"
	secondexamplev1 "acme.corp/pkg/generated/clientset/versioned/typed/secondexample/v1"
	kcpclient "acme.corp/pkg/kcpexisting/clients/clientset/versioned"
	kcpexamplev1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v1"
	fakeexamplev1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v1/fake"
	kcpexamplev1alpha1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v1alpha1"
	fakeexamplev1alpha1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v1alpha1/fake"
	kcpexamplev1beta1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v1beta1"
	fakeexamplev1beta1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v1beta1/fake"
	kcpexamplev2 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v2"
	fakeexamplev2 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example/v2/fake"
	kcpexample3v1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example3/v1"
	fakeexample3v1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/example3/v1/fake"
	kcpexistinginterfacesv1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/existinginterfaces/v1"
	fakeexistinginterfacesv1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/existinginterfaces/v1/fake"
	kcpsecondexamplev1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/secondexample/v1"
	fakesecondexamplev1 "acme.corp/pkg/kcpexisting/clients/clientset/versioned/typed/secondexample/v1/fake"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *ClusterClientset {
	o := kcptesting.NewObjectTracker(clientscheme.Scheme, clientscheme.Codecs.UniversalDecoder())
	o.AddAll(objects...)

	cs := &ClusterClientset{Fake: &kcptesting.Fake{}, tracker: o}
	cs.discovery = &kcpfakediscovery.FakeDiscovery{Fake: cs.Fake, Cluster: logicalcluster.Wildcard}
	cs.AddReactor("*", "*", kcptesting.ObjectReaction(o))
	cs.AddWatchReactor("*", kcptesting.WatchReaction(o))

	return cs
}

var _ kcpclient.ClusterInterface = (*ClusterClientset)(nil)

// ClusterClientset contains the clients for groups.
type ClusterClientset struct {
	*kcptesting.Fake
	discovery *kcpfakediscovery.FakeDiscovery
	tracker   kcptesting.ObjectTracker
}

// Discovery retrieves the DiscoveryClient
func (c *ClusterClientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *ClusterClientset) Tracker() kcptesting.ObjectTracker {
	return c.tracker
}

// Example3V1 retrieves the Example3V1ClusterClient.
func (c *ClusterClientset) Example3V1() kcpexample3v1.Example3V1ClusterInterface {
	return &fakeexample3v1.Example3V1ClusterClient{Fake: c.Fake}
}

// ExampleV1 retrieves the ExampleV1ClusterClient.
func (c *ClusterClientset) ExampleV1() kcpexamplev1.ExampleV1ClusterInterface {
	return &fakeexamplev1.ExampleV1ClusterClient{Fake: c.Fake}
}

// ExampleV1alpha1 retrieves the ExampleV1alpha1ClusterClient.
func (c *ClusterClientset) ExampleV1alpha1() kcpexamplev1alpha1.ExampleV1alpha1ClusterInterface {
	return &fakeexamplev1alpha1.ExampleV1alpha1ClusterClient{Fake: c.Fake}
}

// ExampleV1beta1 retrieves the ExampleV1beta1ClusterClient.
func (c *ClusterClientset) ExampleV1beta1() kcpexamplev1beta1.ExampleV1beta1ClusterInterface {
	return &fakeexamplev1beta1.ExampleV1beta1ClusterClient{Fake: c.Fake}
}

// ExampleV2 retrieves the ExampleV2ClusterClient.
func (c *ClusterClientset) ExampleV2() kcpexamplev2.ExampleV2ClusterInterface {
	return &fakeexamplev2.ExampleV2ClusterClient{Fake: c.Fake}
}

// ExistinginterfacesV1 retrieves the ExistinginterfacesV1ClusterClient.
func (c *ClusterClientset) ExistinginterfacesV1() kcpexistinginterfacesv1.ExistinginterfacesV1ClusterInterface {
	return &fakeexistinginterfacesv1.ExistinginterfacesV1ClusterClient{Fake: c.Fake}
}

// SecondexampleV1 retrieves the SecondexampleV1ClusterClient.
func (c *ClusterClientset) SecondexampleV1() kcpsecondexamplev1.SecondexampleV1ClusterInterface {
	return &fakesecondexamplev1.SecondexampleV1ClusterClient{Fake: c.Fake}
}

// Cluster scopes this clientset to one cluster.
func (c *ClusterClientset) Cluster(cluster logicalcluster.Name) client.Interface {
	if cluster == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}
	return &Clientset{
		Fake:      c.Fake,
		discovery: &kcpfakediscovery.FakeDiscovery{Fake: c.Fake, Cluster: cluster},
		tracker:   c.tracker.Cluster(cluster),
		cluster:   cluster,
	}
}

var _ client.Interface = (*Clientset)(nil)

// Clientset contains the clients for groups.
type Clientset struct {
	*kcptesting.Fake
	discovery *kcpfakediscovery.FakeDiscovery
	tracker   kcptesting.ScopedObjectTracker
	cluster   logicalcluster.Name
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *Clientset) Tracker() kcptesting.ScopedObjectTracker {
	return c.tracker
}

// Example3V1 retrieves the Example3V1Client.
func (c *Clientset) Example3V1() example3v1.Example3V1Interface {
	return &fakeexample3v1.Example3V1Client{Fake: c.Fake, Cluster: c.cluster}
}

// ExampleV1 retrieves the ExampleV1Client.
func (c *Clientset) ExampleV1() examplev1.ExampleV1Interface {
	return &fakeexamplev1.ExampleV1Client{Fake: c.Fake, Cluster: c.cluster}
}

// ExampleV1alpha1 retrieves the ExampleV1alpha1Client.
func (c *Clientset) ExampleV1alpha1() examplev1alpha1.ExampleV1alpha1Interface {
	return &fakeexamplev1alpha1.ExampleV1alpha1Client{Fake: c.Fake, Cluster: c.cluster}
}

// ExampleV1beta1 retrieves the ExampleV1beta1Client.
func (c *Clientset) ExampleV1beta1() examplev1beta1.ExampleV1beta1Interface {
	return &fakeexamplev1beta1.ExampleV1beta1Client{Fake: c.Fake, Cluster: c.cluster}
}

// ExampleV2 retrieves the ExampleV2Client.
func (c *Clientset) ExampleV2() examplev2.ExampleV2Interface {
	return &fakeexamplev2.ExampleV2Client{Fake: c.Fake, Cluster: c.cluster}
}

// ExistinginterfacesV1 retrieves the ExistinginterfacesV1Client.
func (c *Clientset) ExistinginterfacesV1() existinginterfacesv1.ExistinginterfacesV1Interface {
	return &fakeexistinginterfacesv1.ExistinginterfacesV1Client{Fake: c.Fake, Cluster: c.cluster}
}

// SecondexampleV1 retrieves the SecondexampleV1Client.
func (c *Clientset) SecondexampleV1() secondexamplev1.SecondexampleV1Interface {
	return &fakesecondexamplev1.SecondexampleV1Client{Fake: c.Fake, Cluster: c.cluster}
}
