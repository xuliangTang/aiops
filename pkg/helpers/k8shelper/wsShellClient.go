package k8shelper

import (
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"strings"
)

var Upgrader websocket.Upgrader

func init() {
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func HandleCommand(ns, pod, container string, client *kubernetes.Clientset,
	config *rest.Config, command []string) remotecommand.Executor {
	option := &v1.PodExecOptions{
		Container: container,
		Command:   command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}

	req := client.CoreV1().RESTClient().Post().Resource("pods").
		Namespace(ns).
		Name(pod).
		SubResource("exec").
		Param("color", "false").
		VersionedParams(
			option,
			scheme.ParameterCodec,
		)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		panic(err)
	}
	return exec
}

type WsShellClient struct {
	client *websocket.Conn
}

func NewWsShellClient(client *websocket.Conn) *WsShellClient {
	return &WsShellClient{client: client}
}

func (this *WsShellClient) Write(p []byte) (n int, err error) {
	err = this.client.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (this *WsShellClient) Read(p []byte) (n int, err error) {
	_, b, err := this.client.ReadMessage()
	if err != nil {
		return 0, err
	}
	return copy(p, string(b)), nil
}

// 是否是多行文本
func isMultiline(s string) bool {
	return strings.ContainsAny(s, "\n")
}

func clearChar(s string) string {
	s = strings.Replace(s, "`", "", -1)
	return s
}
