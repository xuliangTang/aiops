package qdranthelper

import (
	"context"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

const QdrantAddr = "113.45.67.85:6334"

type QdrantClient struct {
	grpcConn *grpc.ClientConn
}

func NewQdrantClient() *QdrantClient {
	conn, err := grpc.Dial(QdrantAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &QdrantClient{grpcConn: conn}
}

func (qc *QdrantClient) Close() {
	qc.grpcConn.Close()
}

func (qc *QdrantClient) Collection() pb.CollectionsClient {
	return pb.NewCollectionsClient(qc.grpcConn)
}

func toPayload(payload map[string]string) map[string]*pb.Value {
	ret := make(map[string]*pb.Value)
	for k, v := range payload {
		ret[k] = &pb.Value{Kind: &pb.Value_StringValue{StringValue: v}}
	}
	return ret
}

func (qc *QdrantClient) DeleteCollection(name string) error {
	cc := pb.NewCollectionsClient(qc.grpcConn)
	_, err := cc.Delete(context.TODO(), &pb.DeleteCollection{
		CollectionName: name,
	})
	return err
}

// CreateCollection 创建集合
func (qc *QdrantClient) CreateCollection(name string, size uint64) error {
	cc := pb.NewCollectionsClient(qc.grpcConn)
	getReq := pb.GetCollectionInfoRequest{
		CollectionName: name,
	}

	rsp, err := cc.Get(context.Background(), &getReq)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return err
		}
	}

	if rsp.GetResult() != nil {
		log.Println("collection " + name + " already exists,skip create")
		return nil
	}

	req := &pb.CreateCollection{
		CollectionName: name,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     size,
					Distance: pb.Distance_Cosine, // 余弦相似性
				},
			},
		},
	}
	_, err = cc.Create(context.Background(), req)
	if err != nil {
		panic(err)
	}
	return nil
}

// CreatePoints 批量创建Point
func (qc *QdrantClient) CreatePoints(collection string, points []*pb.PointStruct) error {
	pc := pb.NewPointsClient(qc.grpcConn)

	wait := true
	pointsReq := pb.UpsertPoints{
		CollectionName: collection,
		Points:         points,
		Wait:           &wait,
	}

	_, err := pc.Upsert(context.TODO(), &pointsReq)
	if err != nil {
		return err
	}
	return nil
}

// CreatePoint 创建Point的函数
func (qc *QdrantClient) CreatePoint(uuid string, collection string, vector []float32, payload map[string]string) error {
	point := &pb.PointStruct{}
	point.Id = &pb.PointId{
		PointIdOptions: &pb.PointId_Uuid{
			Uuid: uuid,
		},
	}
	point.Vectors = &pb.Vectors{
		VectorsOptions: &pb.Vectors_Vector{
			Vector: &pb.Vector{
				Data: vector,
			},
		},
	}
	point.Payload = toPayload(payload)

	pc := pb.NewPointsClient(qc.grpcConn)

	wait := true
	points := pb.UpsertPoints{
		CollectionName: collection,
		Points:         []*pb.PointStruct{point},
		Wait:           &wait,
	}

	_, err := pc.Upsert(context.TODO(), &points)
	if err != nil {
		return err
	}
	return nil
}

func (qc *QdrantClient) Search(collection string, vector []float32) ([]*pb.ScoredPoint, error) {
	sc := pb.NewPointsClient(qc.grpcConn)
	rsp, err := sc.Search(context.Background(), &pb.SearchPoints{
		CollectionName: collection,
		Vector:         vector,
		Limit:          3, //只取 3条
		WithPayload: &pb.WithPayloadSelector{
			SelectorOptions: &pb.WithPayloadSelector_Include{
				Include: &pb.PayloadIncludeSelector{
					Fields: []string{"url", "method", "body_template"}, // 暴露的字段
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}
