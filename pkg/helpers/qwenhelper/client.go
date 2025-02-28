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

func addHeader(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("llm_qwen")))
	req.Header.Set("Content-Type", "application/json")
}

func doPost(url string, reqBody, respObj interface{}) {
	// 将请求体结构体编码为 JSON
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("Error marshalling request body: %v", err)
	}

	// 构建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// 添加请求头
	addHeader(req)

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
	if err = json.Unmarshal(body, &respObj); err != nil {
		log.Fatalln(err)
	}
}
