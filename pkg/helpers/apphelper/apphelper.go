package apphelper

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"strings"
)

// BuildUrlForGET 用来拼凑url
func BuildUrlForGET(url, bodyTemplate string) (json.RawMessage, error) {
	// 第一步是必须要把body_template 变成JSON对象
	body := make(map[string]string)
	err := json.Unmarshal([]byte(bodyTemplate), &body)
	if err != nil {
		return nil, nil
	}
	for k, v := range body {
		if strings.Index(v, "nothing") >= 0 {
			v = ""
		}
		url = strings.Replace(url, "{"+k+"}", v, -1)
	}

	log.Println("url:", url)
	rsp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, nil
	}
	defer rsp.Body.Close()
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, nil
	}
	getData := gjson.Get(string(b), "data").String()
	return json.RawMessage(getData), nil
}
