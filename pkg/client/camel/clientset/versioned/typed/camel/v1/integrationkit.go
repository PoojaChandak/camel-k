/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	v1 "github.com/apache/camel-k/pkg/apis/camel/v1"
	scheme "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// IntegrationKitsGetter has a method to return a IntegrationKitInterface.
// A group's client should implement this interface.
type IntegrationKitsGetter interface {
	IntegrationKits(namespace string) IntegrationKitInterface
}

// IntegrationKitInterface has methods to work with IntegrationKit resources.
type IntegrationKitInterface interface {
	Create(*v1.IntegrationKit) (*v1.IntegrationKit, error)
	Update(*v1.IntegrationKit) (*v1.IntegrationKit, error)
	UpdateStatus(*v1.IntegrationKit) (*v1.IntegrationKit, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.IntegrationKit, error)
	List(opts metav1.ListOptions) (*v1.IntegrationKitList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.IntegrationKit, err error)
	IntegrationKitExpansion
}

// integrationKits implements IntegrationKitInterface
type integrationKits struct {
	client rest.Interface
	ns     string
}

// newIntegrationKits returns a IntegrationKits
func newIntegrationKits(c *CamelV1Client, namespace string) *integrationKits {
	return &integrationKits{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the integrationKit, and returns the corresponding integrationKit object, and an error if there is any.
func (c *integrationKits) Get(name string, options metav1.GetOptions) (result *v1.IntegrationKit, err error) {
	result = &v1.IntegrationKit{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("integrationkits").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of IntegrationKits that match those selectors.
func (c *integrationKits) List(opts metav1.ListOptions) (result *v1.IntegrationKitList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.IntegrationKitList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("integrationkits").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested integrationKits.
func (c *integrationKits) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("integrationkits").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a integrationKit and creates it.  Returns the server's representation of the integrationKit, and an error, if there is any.
func (c *integrationKits) Create(integrationKit *v1.IntegrationKit) (result *v1.IntegrationKit, err error) {
	result = &v1.IntegrationKit{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("integrationkits").
		Body(integrationKit).
		Do().
		Into(result)
	return
}

// Update takes the representation of a integrationKit and updates it. Returns the server's representation of the integrationKit, and an error, if there is any.
func (c *integrationKits) Update(integrationKit *v1.IntegrationKit) (result *v1.IntegrationKit, err error) {
	result = &v1.IntegrationKit{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("integrationkits").
		Name(integrationKit.Name).
		Body(integrationKit).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *integrationKits) UpdateStatus(integrationKit *v1.IntegrationKit) (result *v1.IntegrationKit, err error) {
	result = &v1.IntegrationKit{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("integrationkits").
		Name(integrationKit.Name).
		SubResource("status").
		Body(integrationKit).
		Do().
		Into(result)
	return
}

// Delete takes name of the integrationKit and deletes it. Returns an error if one occurs.
func (c *integrationKits) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("integrationkits").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *integrationKits) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("integrationkits").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched integrationKit.
func (c *integrationKits) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.IntegrationKit, err error) {
	result = &v1.IntegrationKit{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("integrationkits").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
