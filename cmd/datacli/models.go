package datacli

import (
	"aipos/pkg/helpers/qwenhelper"
	"crypto/md5"
	"encoding/hex"
	pb "github.com/qdrant/go-client/qdrant"
	"log"
	"strings"
)

type PointPayload struct {
	Url          string `json:"url"`
	Method       string `json:"method"`
	BodyTemplate string `json:"body_template"`
}

type Point struct {
	Payload *PointPayload `json:"payload"`
	Prompt  string        `json:"prompt"`
	Ignore  bool          `json:"ignore"` // 如果设置为true则不会进行向量处理
}

func md5str(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func (p *Point) toPayload(host, port string) map[string]*pb.Value {
	ret := make(map[string]*pb.Value)
	// 以下规则 目前是写死的
	p.Payload.Url = strings.Replace(p.Payload.Url, "{host}", host, -1)
	p.Payload.Url = strings.Replace(p.Payload.Url, "{port}", port, -1)
	ret["url"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: p.Payload.Url}}
	ret["method"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: p.Payload.Method}}
	ret["body_template"] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: p.Payload.BodyTemplate}}

	return ret
}

func (p *Point) Build(host, port string) (*pb.PointStruct, error) {
	//vec, err := openaihelper.SimpleGetVec(p.Prompt) //调用OpenAI 获得向量
	vec, err := qwenhelper.GetVec(p.Prompt)
	if err != nil {
		log.Println("获取向量失败", err)
		return nil, err
	}

	// 开始构建PointStruct
	ps := &pb.PointStruct{}

	ps.Id = &pb.PointId{
		PointIdOptions: &pb.PointId_Uuid{
			Uuid: md5str(p.Prompt), // 保证ID唯一
		},
	}
	ps.Vectors = &pb.Vectors{
		VectorsOptions: &pb.Vectors_Vector{
			Vector: &pb.Vector{
				Data: vec,
			},
		},
	}
	ps.Payload = p.toPayload(host, port)
	return ps, nil
}
