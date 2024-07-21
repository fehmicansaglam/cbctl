package couchbase

import (
	"fmt"
)

type DocumentResponse struct {
	Json string `json:"json"`
}

func SearchDocuments(
	bucket string,
	id string,
) (*DocumentResponse, error) {
	endpoint := fmt.Sprintf("pools/default/buckets/%s/docs/%s", bucket, id)
	var response DocumentResponse
	err := getJSONResponse(endpoint, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
