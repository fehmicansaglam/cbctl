package couchbase

import (
	"fmt"
)

type Bucket struct {
	Name         string     `json:"name"`
	BucketType   string     `json:"bucketType"`
	RAMQuotaMB   int        `json:"ramQuotaMB"`
	NumReplicas  int        `json:"numReplicas"`
	FlushEnabled bool       `json:"flushEnabled"`
	BasicStats   BasicStats `json:"basicStats"`
}

type BasicStats struct {
	ItemCount int `json:"itemCount"`
}

func GetBuckets() ([]Bucket, error) {
	var buckets []Bucket

	if err := getJSONResponse("pools/default/buckets", &buckets); err != nil {
		return nil, fmt.Errorf("failed to get index mappings: %w", err)
	}

	return buckets, nil
}
