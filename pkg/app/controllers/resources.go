package controllers

import (
	"aipos/pkg/helpers"
	"aipos/pkg/helpers/k8shelper"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"net/http"
	"strings"
)

// ResourcesController 获取或操作资源的控制器
type ResourcesController struct {
	RestMapper *meta.RESTMapper       `inject:"-"`
	RestConfig *rest.Config           `inject:"-"`
	Client     *dynamic.DynamicClient `inject:"-"`
}

func NewResourcesController() *ResourcesController {
	return &ResourcesController{}
}

// 获取k8s资源
func (ths *ResourcesController) get(ctx *gin.Context) any {
	resName := ctx.Param("resource")
	if resName == "" {
		panic("name is empty")
	}
	mapping, err := k8shelper.MappingFor(resName, ths.RestMapper)
	athena.Error(err)
	ns := ctx.Query("namespace")
	if len(ns) == 0 {
		ns = "default"
	}

	// /resources/get/pods?label[version]=v1&label[app]=nginx
	var labelSelector string
	if labelSet, ok := ctx.GetQueryMap("label"); ok {
		labelSelector = labels.SelectorFromSet(labelSet).String()
	}

	listData, err := k8shelper.ListResource(mapping, ths.Client, ns, labelSelector)
	return listData
}

// 描述k8s资源
func (ths *ResourcesController) describe(ctx *gin.Context) any {
	resName := ctx.Param("resource")
	if resName == "" {
		panic("name is empty")
	}
	objName := ctx.Param("name")
	if objName == "" {
		panic("obj is empty")
	}

	resName = strings.ToLower(helpers.ToPlural(resName))
	mapping, err := k8shelper.MappingFor(resName, ths.RestMapper)
	athena.Error(err)

	ns := ctx.Query("namespace")
	if ns == "" && mapping.Scope.Name() == "namespace" {
		ns = "default"
	}
	ret, err := k8shelper.DescribeResource(mapping, ths.RestConfig, ns, objName)
	athena.Error(err)
	return ret
}

func (ths *ResourcesController) Build(a *athena.Athena) {
	a.Handle(http.MethodGet, "/resources/:resource", ths.get)
	a.Handle(http.MethodGet, "/resources/:resource/:name", ths.describe)
}
