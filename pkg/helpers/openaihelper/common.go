package openaihelper

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"net/url"
	"os"
)

// 设置自己的科学代理地址
func myProxyTransport() *http.Transport {
	SocksProxy := "socks5://127.0.0.1:7890"

	uri, err := url.Parse(SocksProxy)
	if err != nil {
		log.Fatalln(err)
	}
	return &http.Transport{
		Proxy: http.ProxyURL(uri),
	}
}

func NewOpenAiClient() *openai.Client {
	token := os.Getenv("llm_openai")
	config := openai.DefaultConfig(token)
	config.HTTPClient.Transport = myProxyTransport()
	return openai.NewClientWithConfig(config)
}

// SimpleGetVec 把搜索词变成向量
func SimpleGetVec(prompt string) ([]float32, error) {
	c := NewOpenAiClient()
	req := openai.EmbeddingRequest{
		Input: []string{prompt},
		Model: openai.AdaEmbeddingV2,
	}
	rsp, err := c.CreateEmbeddings(context.TODO(), req)
	if err != nil {
		return nil, err
	}
	return rsp.Data[0].Embedding, nil
}
