package controllers

import (
	"aipos/pkg/helpers/k8shelper"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"net/http"
)

// ResourcesController 获取或操作资源的控制器
type ResourcesController struct {
	RestMapper *meta.RESTMapper       `inject:"-"`
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
	ns := ctx.Query("ns")
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

func (ths *ResourcesController) Build(a *athena.Athena) {
	a.Handle(http.MethodGet, "/resources/:resource", ths.get)
}
