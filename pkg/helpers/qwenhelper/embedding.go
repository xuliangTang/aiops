package qwenhelper

import (
	"log"
)

type EmbeddingReq struct {
	Model          string `json:"model"`
	Input          string `json:"input"`
	Dimension      string `json:"dimension"`
	EncodingFormat string `json:"encoding_format"`
}

type EmbeddingResp struct {
	Id     string `json:"id"`
	Object string `json:"object"`
	Model  string `json:"model"`
	Usage  struct {
		PromptTokens int32 `json:"prompt_tokens"`
		TotalTokens  int32 `json:"total_tokens"`
	}
	Data []*struct {
		Embedding []float32
		Index     int32
	}
}

const (
	UrlEmbedding = "https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings"
)

func GetVec(prompt string) ([]float32, error) {
	// 构造请求体
	reqBody := EmbeddingReq{
		Model:          "text-embedding-v3",
		Input:          prompt,
		Dimension:      "1024",
		EncodingFormat: "float",
	}

	vecRet := &EmbeddingResp{}
	doPost(UrlEmbedding, reqBody, vecRet)
	if len(vecRet.Data) == 0 {
		log.Fatalln("empty resp data")
	}

	return vecRet.Data[0].Embedding, nil
}
