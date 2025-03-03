package controllers

import (
	"aipos/pkg/helpers/k8shelper"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuliangTang/athena/athena"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
)

// ShellController 主要处理websocket,包含shell容器、日志查看等
type ShellController struct {
	Client *kubernetes.Clientset `inject:"-"`
	Config *rest.Config          `inject:"-"`
}

func NewShellController() *ShellController {
	return &ShellController{}
}

func (ths *ShellController) firstContainer(pod, ns string) (string, error) {
	// 获取第一个容器 。暂时先这么做  后面再改
	pods, err := ths.Client.CoreV1().Pods(ns).Get(context.TODO(), pod, metav1.GetOptions{})
	athena.Error(err)
	if len(pods.Spec.Containers) > 0 {
		return pods.Spec.Containers[0].Name, nil
	}
	return "", fmt.Errorf("no container found")
}

var ShellCommand = []string{"sh", "-c", "command -v bash >/dev/null 2>&1 && exec bash || exec sh"}

// PodPreConnect 返回一个url给前端，让前端再次连接
func (ths *ShellController) podPreConnect(c *gin.Context) any {
	ns := c.Query("namespace")
	if ns == "" {
		ns = "default"
	}
	pod := c.Query("pod")
	if pod == "" {
		panic("pod name is empty")
	}
	getContainer := c.Query("container")
	if getContainer == "" {
		container, err := ths.firstContainer(pod, ns)
		athena.Error(err)
		getContainer = container
	}

	// 固定shell前缀代表需要客户端重连
	return "shell:/ws/podshell?ns=" + ns + "&pod=" + pod + "&c=" + getContainer
}

func (ths *ShellController) podConnect(c *gin.Context) (v athena.Void) {
	ns := c.Query("namespace")
	if ns == "" {
		ns = "default"
	}
	pod := c.Query("pod")
	getContainer := c.Query("container")
	if getContainer == "" {
		container, err := ths.firstContainer(pod, ns)
		athena.Error(err)
		getContainer = container
	}

	wsClient, err := k8shelper.Upgrader.Upgrade(c.Writer, c.Request, nil)
	athena.Error(err)

	shellClient := k8shelper.NewWsShellClient(wsClient)
	err = k8shelper.HandleCommand(ns, pod, getContainer, ths.Client, ths.Config, ShellCommand).
		StreamWithContext(c, remotecommand.StreamOptions{
			Stdin:  shellClient,
			Stdout: shellClient,
			Stderr: shellClient,
			Tty:    true,
		})

	return
}

func (ths *ShellController) Build(athena *athena.Athena) {
	athena.Handle(http.MethodGet, "/ws/preshell", ths.podPreConnect)
	athena.Handle(http.MethodGet, "/ws/podshell", ths.podConnect)
}
