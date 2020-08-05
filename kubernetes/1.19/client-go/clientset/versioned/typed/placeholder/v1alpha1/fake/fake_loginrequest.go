/*
Copyright 2020 VMware, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/suzerain-io/placeholder-name/kubernetes/1.19/api/apis/placeholder/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeLoginRequests implements LoginRequestInterface
type FakeLoginRequests struct {
	Fake *FakePlaceholderV1alpha1
}

var loginrequestsResource = schema.GroupVersionResource{Group: "placeholder.suzerain-io.github.io", Version: "v1alpha1", Resource: "loginrequests"}

var loginrequestsKind = schema.GroupVersionKind{Group: "placeholder.suzerain-io.github.io", Version: "v1alpha1", Kind: "LoginRequest"}

// Get takes name of the loginRequest, and returns the corresponding loginRequest object, and an error if there is any.
func (c *FakeLoginRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.LoginRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(loginrequestsResource, name), &v1alpha1.LoginRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoginRequest), err
}

// List takes label and field selectors, and returns the list of LoginRequests that match those selectors.
func (c *FakeLoginRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.LoginRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(loginrequestsResource, loginrequestsKind, opts), &v1alpha1.LoginRequestList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.LoginRequestList{ListMeta: obj.(*v1alpha1.LoginRequestList).ListMeta}
	for _, item := range obj.(*v1alpha1.LoginRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested loginRequests.
func (c *FakeLoginRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(loginrequestsResource, opts))
}

// Create takes the representation of a loginRequest and creates it.  Returns the server's representation of the loginRequest, and an error, if there is any.
func (c *FakeLoginRequests) Create(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.CreateOptions) (result *v1alpha1.LoginRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(loginrequestsResource, loginRequest), &v1alpha1.LoginRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoginRequest), err
}

// Update takes the representation of a loginRequest and updates it. Returns the server's representation of the loginRequest, and an error, if there is any.
func (c *FakeLoginRequests) Update(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.UpdateOptions) (result *v1alpha1.LoginRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(loginrequestsResource, loginRequest), &v1alpha1.LoginRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoginRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeLoginRequests) UpdateStatus(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.UpdateOptions) (*v1alpha1.LoginRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(loginrequestsResource, "status", loginRequest), &v1alpha1.LoginRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoginRequest), err
}

// Delete takes name of the loginRequest and deletes it. Returns an error if one occurs.
func (c *FakeLoginRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(loginrequestsResource, name), &v1alpha1.LoginRequest{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLoginRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(loginrequestsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.LoginRequestList{})
	return err
}

// Patch applies the patch and returns the patched loginRequest.
func (c *FakeLoginRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.LoginRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(loginrequestsResource, name, pt, data, subresources...), &v1alpha1.LoginRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LoginRequest), err
}
