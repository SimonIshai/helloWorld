package model

import (
	"fmt"
	"net/http"
)

//type TestCases []TestCase

type TestCase struct {
	ID                     int          `yaml:"id"`
	Randomize              bool         `yaml:"randomize"`
	Request                Request      `yaml:"request"`
	BatchRequest           BatchRequest `yaml:"batch_request"`
	NumOfConcurrentWorkers int          `yaml:"concurrent_workers"`
	NumOfRepetitions       int          `yaml:"repetitions"`
	NumOfRequestsInBatch   int          `yaml:"requests_in_batch"`
	IsBatch                bool         `yaml:"is_batch"`
	client                 *http.Client
}

func (t *TestCase) Log(repetition, workerID int, request, response string) string {

	return fmt.Sprintf(`
CASE ID %d
Repetition %d of %d
WorkerID %d of %d
Is batch - %t
NumOfRequestsInBatch %d

Request:
%s
Response:
%s
------------------------------------------------------------------------------------------------
`, t.ID, repetition, t.NumOfRepetitions, workerID, t.NumOfConcurrentWorkers, t.IsBatch, t.NumOfRequestsInBatch, request, response)

}

func GetTestCases_MaxNumOfRequestsInBatch() []TestCase {
	return []TestCase{

		{ID: 1, NumOfConcurrentWorkers: 100, NumOfRepetitions: 300, NumOfRequestsInBatch: 4, IsBatch: true},
		//{ID: 2, NumOfConcurrentWorkers: 1, NumOfRepetitions: 1, NumOfRequestsInBatch: 2, IsBatch: true},
		//{ID: 3, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 3, IsBatch: true},
		//{ID: 4, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 4, IsBatch: true},
		//{ID: 5, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 5, IsBatch: true},
		//{ID: 6, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 6, IsBatch: true},
		//{ID: 7, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 7, IsBatch: true},
		//{ID: 8, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 8, IsBatch: true},
		//{ID: 9, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 9, IsBatch: true},
		//{ID: 10, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 10, IsBatch: true},
		//{ID: 11, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 11, IsBatch: true},
		//{ID: 12, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 12, IsBatch: true},
		//{ID: 13, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 13, IsBatch: true},
		//{ID: 14, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 14, IsBatch: true},
		//{ID: 15, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 15, IsBatch: true},
		//{ID: 16, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 16, IsBatch: true},
		//{ID: 17, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 17, IsBatch: true},
		//{ID: 18, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 18, IsBatch: true},
		//{ID: 19, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 19, IsBatch: true},
		//{ID: 20, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 20, IsBatch: true},
		//{ID: 21, NumOfConcurrentWorkers: 10, NumOfRepetitions: 1, NumOfRequestsInBatch: 21, IsBatch: true},
	}
}

func GetTestCases_MaxParallelBatches() []TestCase {
	return []TestCase{
		{ID: 1, NumOfConcurrentWorkers: 1, NumOfRepetitions: 120100, NumOfRequestsInBatch: 20, IsBatch: true},
	}
}
