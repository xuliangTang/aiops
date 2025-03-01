package main

import (
	"aipos/pkg/helpers/qdranthelper"
	"aipos/pkg/helpers/qwenhelper"
	"fmt"
	"log"
)

func main() {
	// 提示词生成向量
	userPrompt := "我要获取test命名空间下的所有pods列表" // 描述pod:nginx-54fcc9d5fb-swxgz
	vec, err := qwenhelper.GetVec(userPrompt)
	if err != nil {
		panic(err)
	}

	// qdrant搜素相似的结果，包含api地址和body
	collectionName := "k8smanager"
	points, err := qdranthelper.FastQdrantClient.Search(collectionName, vec)
	if err != nil {
		panic(err)
	}
	ret := points[0]
	if ret.Score <= 0.5 {
		log.Fatalln("暂时不支持你的操作")
	}

	// 让llm替换body_template的内容，生成请求body
	bodyTemplate := ret.Payload["body_template"].GetStringValue()
	fmt.Println(qwenhelper.K8sChat(userPrompt, bodyTemplate))
}
