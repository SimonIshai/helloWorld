package runner

import "fmt"

type TestCases []TestCase

type TestCase struct {
	ID                     int          `json:"id"`
	Request                Request      `json:"request"`
	BatchRequest           BatchRequest `json:"batch_request"`
	NumOfConcurrentWorkers int          `json:"num_of_concurrent_workers"`
	NumOfRepeats           int          `json:"num_of_repeats"`
	IsBatch                bool         `json:"is_batch"`

	RequestMsg  string
	ResponseMsg string
}

var bearerToken = `eyJ0eXAiOiJKV1QiLCJub25jZSI6Ikdjd0NuaW50UW1hUXBjbFVIbDQ3c3FNaHJ2QUtudlBISHRsTXF5bXVSMTgiLCJhbGciOiJSUzI1NiIsIng1dCI6IkJCOENlRlZxeWFHckdOdWVoSklpTDRkZmp6dyIsImtpZCI6IkJCOENlRlZxeWFHckdOdWVoSklpTDRkZmp6dyJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9jM2M5ZGFiZS0xZjg4LTQwOTgtOGIzYS04NTBhNzQyZGI3NmEvIiwiaWF0IjoxNTc0ODY4MDI3LCJuYmYiOjE1NzQ4NjgwMjcsImV4cCI6MTU3NDg3MTkyNywiYWlvIjoiNDJWZ1lLZ29NRXVYZmRXVk5tVk5ST1EvbTlhVkFBPT0iLCJhcHBfZGlzcGxheW5hbWUiOiJhd3NldS1hcG9sbG8tZGV2IiwiYXBwaWQiOiJiMzcyYWEzYy05ZWFjLTQ4MmUtOTRmZi04MzZjNDNlYTIxM2YiLCJhcHBpZGFjciI6IjEiLCJpZHAiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9jM2M5ZGFiZS0xZjg4LTQwOTgtOGIzYS04NTBhNzQyZGI3NmEvIiwib2lkIjoiNTllYzBiZmMtOGJlMi00ZmViLTgyNjMtNmRhNGI1Y2FhNGFhIiwicm9sZXMiOlsiTWFpbC5SZWFkV3JpdGUiLCJHcm91cC5SZWFkLkFsbCIsIlVzZXIuUmVhZC5BbGwiLCJNYWlsLlNlbmQiXSwic3ViIjoiNTllYzBiZmMtOGJlMi00ZmViLTgyNjMtNmRhNGI1Y2FhNGFhIiwidGlkIjoiYzNjOWRhYmUtMWY4OC00MDk4LThiM2EtODUwYTc0MmRiNzZhIiwidXRpIjoiVmpqRkZlbWJ1a3laQjNVdHdYV1NBQSIsInZlciI6IjEuMCIsInhtc190Y2R0IjoxNTcwNDU0MTQ4fQ.WVD0Vygidn_KhipX_lObYPsgOjIixCPDvMCp81ncNe5Y-7uLh6ghJvoDKpsj12kJ3b1_0BVbi-4lOZSREUR5yIvAYTEGQd4-uqp6lqwcI5hcPW0Hby1j7mXMtHbgJNPeR-l1UH_cHl1fWOi1gNu_rvXi4WILV0teNQNF4U_B7tgHFDSHO4y9qXqKbBsxz4j9P49pX34FHj6ssKfmf_6UorjnKem0rMwkulkHwr2IwGJnSWHG9UBD1s8IA9-TMw_KE-nW4VbDHAIDzye0H3ou9Y0k9cF1_VSv8cXMyZ-5G21Cc5FszPLAVGnEirQsjG5Bb3jy2dBl3XQ-Cev8R5v8Tg`

func GetTestCases() TestCases {
	return TestCases{

		//{ID: 1, NumOfConcurrentWorkers: 20, NumOfRepeats: 100, Request: req1, IsBatch: false},
		{ID: 1, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 1), IsBatch: true},
		{ID: 2, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 2), IsBatch: true},
		{ID: 3, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 3), IsBatch: true},
		{ID: 4, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 4), IsBatch: true},
		{ID: 5, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 5), IsBatch: true},
		{ID: 6, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 6), IsBatch: true},
		{ID: 7, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 7), IsBatch: true},
		{ID: 8, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 8), IsBatch: true},
		{ID: 9, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 9), IsBatch: true},
		{ID: 10, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 10), IsBatch: true},
		{ID: 11, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 11), IsBatch: true},
		{ID: 12, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 12), IsBatch: true},
		{ID: 13, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 13), IsBatch: true},
		{ID: 14, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 14), IsBatch: true},
		{ID: 15, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 15), IsBatch: true},
		{ID: 16, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 16), IsBatch: true},
		{ID: 17, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 17), IsBatch: true},
		{ID: 18, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 18), IsBatch: true},
		{ID: 19, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 19), IsBatch: true},
		{ID: 20, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 20), IsBatch: true},
		{ID: 21, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 21), IsBatch: true},
		{ID: 22, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 22), IsBatch: true},
		{ID: 23, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 23), IsBatch: true},
		{ID: 24, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 24), IsBatch: true},
		{ID: 25, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 25), IsBatch: true},
		{ID: 26, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 26), IsBatch: true},
		{ID: 27, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 27), IsBatch: true},
		{ID: 28, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 28), IsBatch: true},
		{ID: 29, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 29), IsBatch: true},
		{ID: 30, NumOfConcurrentWorkers: 10, NumOfRepeats: 1, BatchRequest: generateBatch(batch1.URL, 30), IsBatch: true},
	}
}

func (t *TestCase) String() string {
	return fmt.Sprintf(`
CASE ID %d
Num of repeats %d
Num of concurrent workers %d
Request:
%s
Response:
%s
------------------------------------------------------------------------------------------------
`, t.ID, t.NumOfRepeats, t.NumOfConcurrentWorkers, t.RequestMsg, t.ResponseMsg,
	)
	//return t.RequestMsg + "\n" + t.ResponseMsg
}
