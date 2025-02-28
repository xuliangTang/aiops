package qdranthelper

import (
	"github.com/google/uuid"
)

var FastQdrantClient *QdrantClient

func init() {
	FastQdrantClient = NewQdrantClient()
}

func FastAddPoints(cname string, vecs []float32, payload map[string]string) error {
	uid, _ := uuid.NewUUID()
	return FastQdrantClient.CreatePoint(uid.String(), cname, vecs, payload)
}

type Collection string

func (c Collection) Create(size uint64) error {
	return FastQdrantClient.CreateCollection(string(c), size)
}

func (c Collection) Delete() error {
	return FastQdrantClient.DeleteCollection(string(c))
}
