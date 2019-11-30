package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	gabs "github.com/Jeffail/gabs/v2"
	"github.com/SimonIshai/helloWorld/errors"
	"github.com/SimonIshai/helloWorld/model"
)

var urlGraphAPI = "https://graph.microsoft.com/v1.0/" //users/01242674-b1d8-473d-8558-fe9e043eebcc/"
var bearerToken = `eyJ0eXAiOiJKV1QiLCJub25jZSI6IjRpS2h3MjFySDVHOG51c2V5Tm9JN0VyNnlmb19jcFJneTU2ZlUxY3VsRzQiLCJhbGciOiJSUzI1NiIsIng1dCI6IkJCOENlRlZxeWFHckdOdWVoSklpTDRkZmp6dyIsImtpZCI6IkJCOENlRlZxeWFHckdOdWVoSklpTDRkZmp6dyJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9jM2M5ZGFiZS0xZjg4LTQwOTgtOGIzYS04NTBhNzQyZGI3NmEvIiwiaWF0IjoxNTc0OTU0OTM2LCJuYmYiOjE1NzQ5NTQ5MzYsImV4cCI6MTU3NDk1ODgzNiwiYWlvIjoiNDJWZ1lDaU9GTERTTFBhZWVOemd1b09YOHBsZkFBPT0iLCJhcHBfZGlzcGxheW5hbWUiOiJhd3NldS1hcG9sbG8tZGV2IiwiYXBwaWQiOiJiMzcyYWEzYy05ZWFjLTQ4MmUtOTRmZi04MzZjNDNlYTIxM2YiLCJhcHBpZGFjciI6IjEiLCJpZHAiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9jM2M5ZGFiZS0xZjg4LTQwOTgtOGIzYS04NTBhNzQyZGI3NmEvIiwib2lkIjoiNTllYzBiZmMtOGJlMi00ZmViLTgyNjMtNmRhNGI1Y2FhNGFhIiwicm9sZXMiOlsiTWFpbC5SZWFkV3JpdGUiLCJHcm91cC5SZWFkLkFsbCIsIlVzZXIuUmVhZC5BbGwiLCJNYWlsLlNlbmQiXSwic3ViIjoiNTllYzBiZmMtOGJlMi00ZmViLTgyNjMtNmRhNGI1Y2FhNGFhIiwidGlkIjoiYzNjOWRhYmUtMWY4OC00MDk4LThiM2EtODUwYTc0MmRiNzZhIiwidXRpIjoiOEs3RjZpWXVoa3V3UlV0V1lPRFFBQSIsInZlciI6IjEuMCIsInhtc190Y2R0IjoxNTcwNDU0MTQ4fQ.jWtOJvt7Mu6WkmVqfQq72ZQpxCk1ZPbkCVw4yu2IgBFFSVXYrpX52Ul-FBAavzYDVt1HYplIQ_1iUBZsNtmkIQLn9g5B4B8oKe6cmJ2t6gpbrNtPaConuqRkddVEZIxtLh6INu-z4zIMvzDBEqESU6uXOjWujktspALJBj9jBKXgNc6y8R5vYgmEodYVfRnQtIWL0SOl88jakXmybxXz4Bcw_Q00RAAm147wAw7u5cZWp9QnLbgaaRHvOYdvxhGRZKPmeSNnrQYsWmgF36djhvzT2TwFXDkZkJ61mOczmzeNlhq0p8dfNaD_WIR-A-j3lfDuh_ANzeFAZs7AvphD1w`
var client = &http.Client{Timeout: 60 * time.Second}

func Run(testCases []model.TestCase) error {

	testCasesNum := len(testCases)
	var wg sync.WaitGroup
	wg.Add(testCasesNum)
	for i := range testCases {

		go func(caseIndex int) {
			defer wg.Done()

			err := runTestCase(&testCases[caseIndex])
			if err != nil {
				log.Printf("TestCase %d, err: %s\n", testCases[caseIndex].ID, err)
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func runTestCase(testCase *model.TestCase) error {

	chanSize := testCase.NumOfRepetitions
	ch := make(chan int, chanSize+1)
	for i := 1; i <= chanSize; i++ {
		ch <- i
	}
	close(ch)
	log.Println("run CaseID", testCase.ID, "NumOfConcurrentWorkers", testCase.NumOfConcurrentWorkers)
	var wg sync.WaitGroup
	wg.Add(testCase.NumOfConcurrentWorkers)
	for i := 1; i <= testCase.NumOfConcurrentWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()

			for repetition := range ch {

				log.Println("CaseID", testCase.ID, "worker", workerID, ", repetition", repetition)

				err := runRepetition(testCase, workerID, repetition)
				if err != nil {
					log.Println(errors.Wrap(errors.Wrap(err, "runRepetition"), "case run").Error())
					return
				}
			}

		}(i)
	}
	wg.Wait()
	return nil
}

func runRepetition(testCase *model.TestCase, workerID, repetition int) error {

	var url string
	var method string
	var body []byte
	var err error

	if testCase.IsBatch {
		testCase.BatchRequest = model.GenerateBatch(testCase.NumOfRequestsInBatch)
		method = testCase.BatchRequest.Method
		url = urlGraphAPI + testCase.BatchRequest.URL
		body, err = json.MarshalIndent(testCase.BatchRequest.Body, "", "\t")
		if err != nil {
			return errors.WrapWithKind(err, errors.KindParse, "json marshal")
		}

	} else {
		url = urlGraphAPI + testCase.Request.URL
		method = testCase.Request.Method
		if testCase.Request.Body != "" {
			body = []byte(testCase.Request.Body)
		}
	}
	requestMsg := method + ": " + url + "\n" + string(body) + "\n"

	respBody, err := sendRequest(method, url, body)
	if err != nil {
		return errors.Wrap(err, "send request to graph api")
	}
	responseMsg := string(respBody)

	jsonParsed, err := gabs.ParseJSON(respBody)
	if err != nil {
		return errors.WrapWithKind(err, errors.KindParse, "parse response to json")
	}

	timestamp := time.Now().Format("2006-01-02T15:04")

	errPath := "error"
	if jsonParsed.Exists(errPath) {

		filename := fmt.Sprintf("%s_CaseID_%d_repetition_%d_worker_%d_error.json", timestamp, testCase.ID, repetition, workerID)
		if err := storeToFile([]byte(testCase.Log(repetition, workerID, requestMsg, responseMsg)), filename); err != nil {
			return errors.Wrap(err, "store response to file")
		}
		return errors.New(errors.KindGraphAPI, "got error message")
	}

	responsesPath := "responses"
	if jsonParsed.Exists(responsesPath) {

		for _, child := range jsonParsed.S(responsesPath).Children() {

			if child.ExistsP("body.error") {

				filename := fmt.Sprintf("%s_CaseID_%d_repetition_%d_worker_%d_error.json", timestamp, testCase.ID, repetition, workerID)
				if err := storeToFile([]byte(testCase.Log(repetition, workerID, requestMsg, responseMsg)), filename); err != nil {
					return errors.Wrap(err, "store response to file")
				}
				return errors.New(errors.KindGraphAPI, "got error message")
			}
		}
	}

	return nil
}

func sendRequest(method, url string, body []byte) ([]byte, error) {

	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, errors.WrapWithKind(err, errors.KindHttp, "NewRequest")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WrapWithKind(err, errors.KindHttp, "client.Do")
	}
	defer resp.Body.Close() //nolint: errcheck

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WrapWithKind(err, errors.KindHttp, "read response body")
	}
	return respBody, nil
}

func storeToFile(bytesToStore []byte, filename string) error {

	logFolder := "test-case-logs"
	_, err := os.Stat(logFolder)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(logFolder, os.ModePerm)
		if errDir != nil {
			return errors.WrapWithKind(err, errors.KindFileSystem, "create logs folder")
		}
	}

	filename = logFolder + "/" + filename
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.WrapWithKind(err, errors.KindFileSystem, "open file")
	}

	defer f.Close()
	bytesToStore = append(bytesToStore, []byte("\n")...)
	if _, err = f.Write(bytesToStore); err != nil {
		return errors.WrapWithKind(err, errors.KindFileSystem, "file write")
	}
	return nil
}
