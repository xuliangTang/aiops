package openaihelper

import (
	"context"
	"encoding/json"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"log"

	"os"
	"strings"
)

var MessageStore ChatMessages

func init() {
	MessageStore = make(ChatMessages, 0)
	MessageStore.Clear() //清理和初始化
}

const PromptTemplate = `please help me extract the k8s elements of this text: "{prompt}"
Here is the answer template:
{body_template}
Please fill in the text in {} completely according to the template. If it cannot be extracted, fill in "nothing". Each line cannot be omitted and must be filled in English. Don't give any explanation`

func K8sChat(userPrompt, bodyTemplate string) string {
	prompt := strings.Replace(PromptTemplate, "{prompt}", userPrompt, 1)
	prompt = strings.Replace(prompt, "{body_template}", bodyTemplate, 1)
	c := NewOpenAiClient()
	MessageStore.AddForUser(prompt)
	rsp, err := c.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: MessageStore.ToMessage(),
	})
	if err != nil {
		log.Println(err)
		return ""
	}
	MessageStore.AddForAssistant(rsp.Choices[0].Message.Content)

	return MessageStore.GetLast()
}

type ChatMessages []*ChatMessage
type ChatMessage struct {
	Msg openai.ChatCompletionMessage
}

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)

func (cm *ChatMessages) Clear() {
	*cm = make([]*ChatMessage, 0) //重新初始化
	msg := `You are a helpful k8s assistant`
	cm.AddForSystem(msg)
}

func (cm *ChatMessages) AddFor(msg string, role string) {
	*cm = append(*cm, &ChatMessage{
		Msg: openai.ChatCompletionMessage{
			Role:    role,
			Content: msg,
		},
	})
}

const CommandPattern = "```\\s*(.*?)\\s*```"

func (cm *ChatMessages) Dump(file string) error {
	//把内容json化后后 ，存到文件里
	b, err := json.Marshal(cm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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

func (cm *ChatMessages) Apply(content string) error {
	content = strings.Trim(content, " ")
	if len(*cm) == 0 || content == "" {
		return fmt.Errorf("无需apply")
	}
	msg := (*cm)[len(*cm)-1]
	if msg.Msg.Role != RoleAssistant {
		fmt.Println("user/system内容无需apply")
		return fmt.Errorf("user/system内容无需apply")
	}
	msg.Msg.Content = content
	return nil
}

func (cm *ChatMessages) ToMessage() []openai.ChatCompletionMessage {
	ret := make([]openai.ChatCompletionMessage, len(*cm))
	for index, c := range *cm {
		ret[index] = c.Msg
	}
	return ret
}

func (cm *ChatMessages) GetLast() string {
	if len(*cm) == 0 {
		return "什么都没找到"
	}

	return (*cm)[len(*cm)-1].Msg.Content
}
