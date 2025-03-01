package k8shelper

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/describe"
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

func DescribeResource(mapping *meta.RESTMapping, restConfig *rest.Config, ns, name string) (string, error) {
	resDescriber, ok := describe.DescriberFor(mapping.GroupVersionKind.GroupKind(), restConfig)
	if !ok {
		return "", fmt.Errorf("resource describe error")
	}

	ret, err := resDescriber.Describe(ns, name, describe.DescriberSettings{})
	if err != nil {
		return "", err
	}
	return ret, nil
}
