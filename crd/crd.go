package crd

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

type Client struct {
	Client    *dynamic.DynamicClient
	Namespace string
	GVR       schema.GroupVersionResource
}

func (c *Client) Get(ctx context.Context, name string, opts metav1.GetOptions) (*unstructured.Unstructured, error) {
	return c.Client.Resource(c.GVR).Namespace(c.Namespace).Get(ctx, name, opts)
}

func (c *Client) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return c.Client.Resource(c.GVR).Namespace(c.Namespace).List(ctx, opts)
}

func (c *Client) Create(ctx context.Context, data interface{}, opts metav1.CreateOptions) (*unstructured.Unstructured, error) {
	obj, ok := data.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("can not parse data =%v to *unstructured.Unstructured", data)
	}

	return c.Client.Resource(c.GVR).Namespace(c.Namespace).Create(ctx, obj, opts)
}

func (c *Client) Update(ctx context.Context, data interface{}, opts metav1.UpdateOptions) (
	*unstructured.Unstructured, error) {
	obj, ok := data.(*unstructured.Unstructured)
	if !ok {
		return nil, fmt.Errorf("can not parse data =%v to *unstructured.Unstructured obj", data)
	}
	return c.Client.Resource(c.GVR).Update(ctx, obj, opts)
}

func (c *Client) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Client.Resource(c.GVR).Delete(ctx, name, opts)
}

func (c *Client) Patch(ctx context.Context, name string, pt types.PatchType, date []byte, opts metav1.PatchOptions) (
	*unstructured.Unstructured,
	error,
) {
	return c.Client.Resource(c.GVR).Namespace(c.Namespace).Patch(ctx, name, pt, date, opts)
}

func (c *Client) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.Client.Resource(c.GVR).Namespace(c.Namespace).Watch(ctx, opts)
}
