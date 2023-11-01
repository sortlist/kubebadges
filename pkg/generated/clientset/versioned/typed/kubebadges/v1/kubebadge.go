/*
Copyright 2023 Your Company.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/kubebadges/kubebadges/pkg/apis/kubebadges/v1"
	scheme "github.com/kubebadges/kubebadges/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// KubeBadgesGetter has a method to return a KubeBadgeInterface.
// A group's client should implement this interface.
type KubeBadgesGetter interface {
	KubeBadges(namespace string) KubeBadgeInterface
}

// KubeBadgeInterface has methods to work with KubeBadge resources.
type KubeBadgeInterface interface {
	Create(ctx context.Context, kubeBadge *v1.KubeBadge, opts metav1.CreateOptions) (*v1.KubeBadge, error)
	Update(ctx context.Context, kubeBadge *v1.KubeBadge, opts metav1.UpdateOptions) (*v1.KubeBadge, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.KubeBadge, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.KubeBadgeList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.KubeBadge, err error)
	KubeBadgeExpansion
}

// kubeBadges implements KubeBadgeInterface
type kubeBadges struct {
	client rest.Interface
	ns     string
}

// newKubeBadges returns a KubeBadges
func newKubeBadges(c *KubebadgesV1Client, namespace string) *kubeBadges {
	return &kubeBadges{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the kubeBadge, and returns the corresponding kubeBadge object, and an error if there is any.
func (c *kubeBadges) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.KubeBadge, err error) {
	result = &v1.KubeBadge{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kubebadges").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KubeBadges that match those selectors.
func (c *kubeBadges) List(ctx context.Context, opts metav1.ListOptions) (result *v1.KubeBadgeList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.KubeBadgeList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("kubebadges").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kubeBadges.
func (c *kubeBadges) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("kubebadges").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a kubeBadge and creates it.  Returns the server's representation of the kubeBadge, and an error, if there is any.
func (c *kubeBadges) Create(ctx context.Context, kubeBadge *v1.KubeBadge, opts metav1.CreateOptions) (result *v1.KubeBadge, err error) {
	result = &v1.KubeBadge{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("kubebadges").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kubeBadge).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a kubeBadge and updates it. Returns the server's representation of the kubeBadge, and an error, if there is any.
func (c *kubeBadges) Update(ctx context.Context, kubeBadge *v1.KubeBadge, opts metav1.UpdateOptions) (result *v1.KubeBadge, err error) {
	result = &v1.KubeBadge{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("kubebadges").
		Name(kubeBadge.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kubeBadge).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the kubeBadge and deletes it. Returns an error if one occurs.
func (c *kubeBadges) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kubebadges").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kubeBadges) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("kubebadges").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched kubeBadge.
func (c *kubeBadges) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.KubeBadge, err error) {
	result = &v1.KubeBadge{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("kubebadges").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
