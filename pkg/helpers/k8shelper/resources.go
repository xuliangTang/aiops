package k8shelper

import (
	"context"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

func ListResource(mapping *meta.RESTMapping, client *dynamic.DynamicClient, ns, label string) (*unstructured.UnstructuredList, error) {
	var ri dynamic.ResourceInterface
	if mapping.Scope.Name() == "namespace" {
		ri = client.Resource(mapping.Resource).Namespace(ns)
	} else {
		ri = client.Resource(mapping.Resource)
	}
	return ri.List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
		Limit:         10,
	})
}
