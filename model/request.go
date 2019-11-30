package model

import (
	"github.com/gofrs/uuid"
)

type Request struct {
	ID      string            `yaml:"id"`
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Body    string            `yaml:"-"`
	Headers map[string]string `yaml:"-"`
}

type BatchRequestBody struct {
	Requests []Request `json:"requests"`
}

type BatchRequest struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   BatchRequestBody
}

var urlRouteBatch = "$batch"

var req1 = Request{
	ID:     "001",
	Method: "GET",
	URL:    "users/user1@sishaicyren.onmicrosoft.com/messages/AAMkAGY2NjI5NzU0LTQxYjktNGQ5My1hODY0LWYxNmM1Zjg1ZDNmZQBGAAAAAAANsx3zYZWoRbgJTxzNLUUrBwBKgViXObVQRKVrwdIXxyCeAAAAAAEMAABKgViXObVQRKVrwdIXxyCeAAARuI43AAA\u003d",
}

var req2 = Request{
	ID:     "002",
	Method: "GET",
	URL:    "users/eb16a110-2c59-4dc8-b2b5-4f24f93c86ba/",
}

var req3 = Request{
	ID:     "003",
	Method: "GET",
	URL:    "users/bf10c02b-ad69-47ad-a217-2a98efd60f79/",
}

var batch1 = BatchRequest{
	Method: "POST",
	URL:    "$batch",
	Body:   BatchRequestBody{[]Request{req1, req2}},
}

func GenerateBatch(numOfRequestsInBatch int) BatchRequest {

	var requests []Request
	for i := 0; i < numOfRequestsInBatch; i++ {
		uuID, _ := uuid.NewV4()
		req := req1
		req.ID = uuID.String()
		requests = append(requests, req)
	}

	return BatchRequest{
		Method: "POST",
		URL:    urlRouteBatch,
		Body:   BatchRequestBody{requests},
	}
}
