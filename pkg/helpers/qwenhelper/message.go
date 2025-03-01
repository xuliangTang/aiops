package qwenhelper

import (
	"log"
	"strings"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

var MessageStore ChatMessages

func init() {
	MessageStore = make(ChatMessages, 0)
	MessageStore.Clear() //清理和初始化
}

type ChatMessages []*ChatCompletionMessage
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (cm *ChatMessages) AddForAssistant(msg string) {
	cm.AddFor(msg, RoleAssistant)
}

func (cm *ChatMessages) AddForSystem(msg string) {
	cm.AddFor(msg, RoleSystem)
}

func (cm *ChatMessages) AddForUser(msg string) {
	cm.AddFor(msg, RoleUser)
}

func (cm *ChatMessages) AddFor(msg string, role string) {
	*cm = append(*cm, &ChatCompletionMessage{
		Role:    role,
		Content: msg,
	})
}

func (cm *ChatMessages) Clear() {
	*cm = make([]*ChatCompletionMessage, 0) //重新初始化
	msg := `You are a helpful k8s assistant`
	cm.AddForSystem(msg)
}

const UrlChat = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

const PromptTemplate = `Extract the k8s elements from this text: "{prompt}"
Use this template:
{body_template}
Fill the {} with extracted values. If a field cannot be extracted, use "nothing". For "name", extract the resource type (e.g., "pods", "services"). Provide the result as plain JSON without markdown or explanations. Each line must be filled in English.`

func K8sChat(userPrompt, bodyTemplate string) string {
	prompt := strings.Replace(PromptTemplate, "{prompt}", userPrompt, 1)
	prompt = strings.Replace(prompt, "{body_template}", bodyTemplate, 1)

	MessageStore.AddForUser(prompt)
	req := &ChatReq{
		Model:    "qwen-max-2025-01-25",
		Messages: MessageStore,
	}
	resp := &ChatResp{}
	doPost(UrlChat, req, resp)
	if len(resp.Choices) == 0 {
		log.Fatalln("empty resp data")
	}

	return resp.Choices[0].Message.Content
}

type ChatReq struct {
	Model    string                   `json:"model"`
	Messages []*ChatCompletionMessage `json:"messages"`
}

type ChatResp struct {
	Id     string `json:"id"`
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  *struct {
		PromptTokens     int32 `json:"prompt_tokens"`
		CompletionTokens int32 `json:"completion_tokens"`
		TotalTokens      int32 `json:"total_tokens"`
	}
	Choices []*struct {
		Index        int32  `json:"index"`
		FinishReason string `json:"finish_reason"`
		Message      *struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}
	}
	Created int64 `json:"created"`
}
