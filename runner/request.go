package runner

import "fmt"

type Request struct {
	ID      string            `json:"id"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Body    string            `json:"body,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type BatchRequestBody struct {
	Requests []Request `json:"requests"`
}

type BatchRequest struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   BatchRequestBody
}

var urlAPI = "https://graph.microsoft.com/v1.0/" //users/01242674-b1d8-473d-8558-fe9e043eebcc/"

var req1 = Request{
	ID:     "001",
	Method: "GET",
	URL:    "users/01242674-b1d8-473d-8558-fe9e043eebcc/",
}

var req2 = Request{
	ID:     "002",
	Method: "GET",
	URL:    "users/01242674-b1d8-473d-8558-fe9e043eebcc/",
}

var batch1 = BatchRequest{
	Method: "POST",
	URL:    "$batch",
	Body:   BatchRequestBody{[]Request{req1, req2}},
}

func generateBatch(url string, numOfRequests int) BatchRequest {

	var requests []Request
	for i := 0; i < numOfRequests; i++ {
		req := req1
		req.ID = fmt.Sprintf("%d", i+1)
		requests = append(requests, req)
	}

	return BatchRequest{
		Method: "POST",
		URL:    url,
		Body:   BatchRequestBody{requests},
	}
}
