package controllers

import (
	"aipos/pkg/helpers/apphelper"
	"aipos/pkg/helpers/qdranthelper"
	"aipos/pkg/helpers/qwenhelper"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena"
	"log"
	"net/http"
)

// PromptController 前端请求的专用控制器
type PromptController struct{}

func NewPromptController() *PromptController {
	return &PromptController{}
}

type PromptRequest struct {
	Prompt string `json:"prompt" form:"prompt" binding:"required"`
}

const collectionName = "k8smanager"

func (ths *PromptController) Do(c *gin.Context) any {
	qwenhelper.MessageStore.Clear()
	req := &PromptRequest{}
	err := c.ShouldBindJSON(req)
	athena.Error(err)

	// 转为向量
	vec, err := qwenhelper.GetVec(req.Prompt)
	athena.Error(err)
	// 从向量数据库中执行相似性搜索
	points, err := qdranthelper.FastQdrantClient.Search(collectionName, vec)
	athena.Error(err)
	ret := points[0]
	if ret.Score < 0.5 {
		panic("暂时不支持你的操作")
	}

	// 得到api，构建请求
	method := ret.Payload["method"].GetStringValue()
	url := ret.Payload["url"].GetStringValue()
	bodyTemplate := ret.Payload["body_template"].GetStringValue()
	body := qwenhelper.K8sChat(req.Prompt, bodyTemplate)
	if body == "" {
		panic("调用通义千问接口失败")
	}
	log.Println("url:", url)
	log.Println("body", body)
	if method == "GET" {
		ret, err := apphelper.BuildUrlForGET(url, body)
		athena.Error(err)
		return ret
	}

	return nil
}

func (ths *PromptController) Build(athena *athena.Athena) {
	athena.Handle(http.MethodPost, "/prompt", ths.Do)
}
