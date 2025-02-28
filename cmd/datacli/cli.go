package datacli

import (
	"aipos/pkg/helpers/qdranthelper"
	"bufio"
	"bytes"
	"fmt"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/apimachinery/pkg/util/yaml"
	"log"
	"os"
)

// LoadPoints 模仿kubectl读取yaml文件
func LoadPoints(file string) ([]*Point, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var ret []*Point

	bf := bufio.NewReader(f)
	// 读取一个文件中多个yaml 并转为指定对象的代码
	r := yaml.NewYAMLReader(bf)
	for {
		b, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		tmp := &Point{}
		d := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(b), 1024*10)
		err = d.Decode(tmp)
		if err != nil {
			return nil, err
		}
		ret = append(ret, tmp)
	}

	return ret, nil
}

const (
	HOST = "localhost"
	PORT = "8080"
)

func Run(f string, collection string) error {
	points, err := LoadPoints(f)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = qdranthelper.Collection(collection).Create(1536)
	if err != nil {
		log.Fatalln("创建collection出错:", err.Error())
	}
	var psSet []*pb.PointStruct
	for _, p := range points {
		if p.Ignore {
			continue
		}
		ps, err := p.Build(HOST, PORT)
		if err != nil {
			log.Println(err)
			continue
		}
		psSet = append(psSet, ps)

	}
	err = qdranthelper.FastQdrantClient.CreatePoints(collection, psSet)
	if err != nil {
		log.Fatalln("批量创建point出错:", err.Error())
	}
	fmt.Println("导入成功")
	return nil
}

// NewDataCliCommand 初始化命令行参数
func NewDataCliCommand() *cobra.Command {
	// 集成 cobra 命令
	cmd := &cobra.Command{
		Use: "dataclient",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			f, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}
			if f == "" {
				return fmt.Errorf("file is empty")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// 获取flag
			f, _ := cmd.Flags().GetString("file")
			c, _ := cmd.Flags().GetString("collection")
			return Run(f, c)
		},
	}
	// 添加flag, name=port,默认值是8080
	cmd.Flags().StringP("file", "f", "", "dataclient -f ")
	cmd.Flags().StringP("collection", "c", "k8smanager", "dataclient -c xxx ")

	return cmd
}
