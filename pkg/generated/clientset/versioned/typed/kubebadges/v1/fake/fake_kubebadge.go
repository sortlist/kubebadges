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

package fake

import (
	"context"

	v1 "github.com/kubebadges/kubebadges/pkg/apis/kubebadges/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKubeBadges implements KubeBadgeInterface
type FakeKubeBadges struct {
	Fake *FakeKubebadgesV1
	ns   string
}

var kubebadgesResource = v1.SchemeGroupVersion.WithResource("kubebadges")

var kubebadgesKind = v1.SchemeGroupVersion.WithKind("KubeBadge")

// Get takes name of the kubeBadge, and returns the corresponding kubeBadge object, and an error if there is any.
func (c *FakeKubeBadges) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.KubeBadge, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(kubebadgesResource, c.ns, name), &v1.KubeBadge{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.KubeBadge), err
}

// List takes label and field selectors, and returns the list of KubeBadges that match those selectors.
func (c *FakeKubeBadges) List(ctx context.Context, opts metav1.ListOptions) (result *v1.KubeBadgeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(kubebadgesResource, kubebadgesKind, c.ns, opts), &v1.KubeBadgeList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.KubeBadgeList{ListMeta: obj.(*v1.KubeBadgeList).ListMeta}
	for _, item := range obj.(*v1.KubeBadgeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested kubeBadges.
func (c *FakeKubeBadges) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(kubebadgesResource, c.ns, opts))

}

// Create takes the representation of a kubeBadge and creates it.  Returns the server's representation of the kubeBadge, and an error, if there is any.
func (c *FakeKubeBadges) Create(ctx context.Context, kubeBadge *v1.KubeBadge, opts metav1.CreateOptions) (result *v1.KubeBadge, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(kubebadgesResource, c.ns, kubeBadge), &v1.KubeBadge{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.KubeBadge), err
}

// Update takes the representation of a kubeBadge and updates it. Returns the server's representation of the kubeBadge, and an error, if there is any.
func (c *FakeKubeBadges) Update(ctx context.Context, kubeBadge *v1.KubeBadge, opts metav1.UpdateOptions) (result *v1.KubeBadge, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(kubebadgesResource, c.ns, kubeBadge), &v1.KubeBadge{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.KubeBadge), err
}

// Delete takes name of the kubeBadge and deletes it. Returns an error if one occurs.
func (c *FakeKubeBadges) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(kubebadgesResource, c.ns, name, opts), &v1.KubeBadge{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKubeBadges) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(kubebadgesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.KubeBadgeList{})
	return err
}

// Patch applies the patch and returns the patched kubeBadge.
func (c *FakeKubeBadges) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.KubeBadge, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(kubebadgesResource, c.ns, name, pt, data, subresources...), &v1.KubeBadge{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.KubeBadge), err
}