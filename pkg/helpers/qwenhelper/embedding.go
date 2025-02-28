package qwenhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type requestBody struct {
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
	reqBody := requestBody{
		Model:          "text-embedding-v3",
		Input:          prompt,
		Dimension:      "1024",
		EncodingFormat: "float",
	}

	// 将请求体结构体编码为 JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("Error marshalling request body: %v", err)
	}

	// 构建请求
	req, err := http.NewRequest("POST", UrlEmbedding, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// 添加请求头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("llm_qwen")))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// 打印响应
	fmt.Printf("Response Status: %s\n", resp.Status)
	vecRet := &EmbeddingResp{}
	if err = json.Unmarshal(body, vecRet); err != nil {
		log.Fatalln(err)
	}
	if len(vecRet.Data) == 0 {
		log.Fatalln("empty resp data")
	}

	return vecRet.Data[0].Embedding, nil
}
