package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	gabs "github.com/Jeffail/gabs/v2"
	"github.com/SimonIshai/helloWorld/config"
	"github.com/SimonIshai/helloWorld/errors"
	"github.com/SimonIshai/helloWorld/model"
)

var urlGraphAPI = `https://graph.microsoft.com/v1.0/`
var bearerToken = `eyJ0eXAiOiJKV1QiLCJub25jZSI6IjRpS2h3MjFySDVHOG51c2V5Tm9JN0VyNnlmb19jcFJneTU2ZlUxY3VsRzQiLCJhbGciOiJSUzI1NiIsIng1dCI6IkJCOENlRlZxeWFHckdOdWVoSklpTDRkZmp6dyIsImtpZCI6IkJCOENlRlZxeWFHckdOdWVoSklpTDRkZmp6dyJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9jM2M5ZGFiZS0xZjg4LTQwOTgtOGIzYS04NTBhNzQyZGI3NmEvIiwiaWF0IjoxNTc0OTU0OTM2LCJuYmYiOjE1NzQ5NTQ5MzYsImV4cCI6MTU3NDk1ODgzNiwiYWlvIjoiNDJWZ1lDaU9GTERTTFBhZWVOemd1b09YOHBsZkFBPT0iLCJhcHBfZGlzcGxheW5hbWUiOiJhd3NldS1hcG9sbG8tZGV2IiwiYXBwaWQiOiJiMzcyYWEzYy05ZWFjLTQ4MmUtOTRmZi04MzZjNDNlYTIxM2YiLCJhcHBpZGFjciI6IjEiLCJpZHAiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9jM2M5ZGFiZS0xZjg4LTQwOTgtOGIzYS04NTBhNzQyZGI3NmEvIiwib2lkIjoiNTllYzBiZmMtOGJlMi00ZmViLTgyNjMtNmRhNGI1Y2FhNGFhIiwicm9sZXMiOlsiTWFpbC5SZWFkV3JpdGUiLCJHcm91cC5SZWFkLkFsbCIsIlVzZXIuUmVhZC5BbGwiLCJNYWlsLlNlbmQiXSwic3ViIjoiNTllYzBiZmMtOGJlMi00ZmViLTgyNjMtNmRhNGI1Y2FhNGFhIiwidGlkIjoiYzNjOWRhYmUtMWY4OC00MDk4LThiM2EtODUwYTc0MmRiNzZhIiwidXRpIjoiOEs3RjZpWXVoa3V3UlV0V1lPRFFBQSIsInZlciI6IjEuMCIsInhtc190Y2R0IjoxNTcwNDU0MTQ4fQ.jWtOJvt7Mu6WkmVqfQq72ZQpxCk1ZPbkCVw4yu2IgBFFSVXYrpX52Ul-FBAavzYDVt1HYplIQ_1iUBZsNtmkIQLn9g5B4B8oKe6cmJ2t6gpbrNtPaConuqRkddVEZIxtLh6INu-z4zIMvzDBEqESU6uXOjWujktspALJBj9jBKXgNc6y8R5vYgmEodYVfRnQtIWL0SOl88jakXmybxXz4Bcw_Q00RAAm147wAw7u5cZWp9QnLbgaaRHvOYdvxhGRZKPmeSNnrQYsWmgF36djhvzT2TwFXDkZkJ61mOczmzeNlhq0p8dfNaD_WIR-A-j3lfDuh_ANzeFAZs7AvphD1w`
var client = &http.Client{Timeout: 60 * time.Second}

func Run(testCases []model.TestCase) error {

	if err := getGraphApiToken(); err != nil {
		return errors.Wrap(err, "getGraphApiToken")
	}

	for i := range testCases {

		func(caseIndex int) {

			err := runTestCase(&testCases[caseIndex])
			if err != nil {
				log.Printf("TestCase %d, err: %s\n", testCases[caseIndex].ID, err)
			}
		}(i)
	}

	log.Println("--------------------------- Statistics Results ---------------------------")
	for i := range testCases {
		//fmt.Println("Case", testCases[i].ID, "Statistics", testCases[i].Statistics.String())
		testCases[i].Statistics.Refresh()
		//statBytes, _ := json.MarshalIndent(&testCases[i].Statistics, "", "  ")
		caseBytes, _ := json.MarshalIndent(&testCases[i], "", "  ")
		//caseBytes, _ := json.Marshal(&testCases[i])
		//fmt.Printf("Case %d \n%s\n", testCases[i].ID, string(statBytes))
		fmt.Printf("%s\n", string(caseBytes))
	}
	return nil
}

func runTestCase(testCase *model.TestCase) error {

	testCase.Statistics.SpanStartTime = time.Now()
	defer func() {
		testCase.Statistics.SpanEndTime = time.Now()
	}()

	if testCase.IsBatch {
		testCase.Statistics.TotalRequests = testCase.NumOfRepetitions * testCase.NumOfRequestsInBatch
	} else {
		testCase.Statistics.TotalRequests = testCase.NumOfRepetitions
	}

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
			defer func() {
				wg.Done()
				log.Println("finishing CaseID", testCase.ID, "worker", workerID)
			}()

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

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", bearerToken),
	}

	respBody, err := sendRequest(method, url, body, headers)
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
		msResponse := model.MSResponse{}
		if err := json.Unmarshal(jsonParsed.Path(errPath).Bytes(), &msResponse.Body.MSError); err != nil {
			return errors.WrapWithKind(err, errors.KindGraphAPI, "store response to file")
		}

		if err := testCase.Statistics.AddMSErrorResponse(&msResponse, testCase.NumOfRequestsInBatch); err != nil {
			return errors.Wrap(err, "AddMSErrorResponse")
		}
		filename := fmt.Sprintf("%s_CaseID_%d_repetition_%d_worker_%d_error.json", timestamp, testCase.ID, repetition, workerID)
		if err := storeToFile([]byte(testCase.Log(repetition, workerID, requestMsg, responseMsg)), filename); err != nil {
			return errors.Wrap(err, "store response to file")
		}
		log.Printf("CaseID_%d_repetition_%d_worker_%d, err: %s", testCase.ID, repetition, workerID, jsonParsed.Path(errPath).String())
	}

	responsesPath := "responses"
	if jsonParsed.Exists(responsesPath) {

		for _, child := range jsonParsed.S(responsesPath).Children() {

			if child.ExistsP("body.error") {
				msResponse := model.MSResponse{}
				if err := json.Unmarshal(child.Bytes(), &msResponse); err != nil {
					return errors.WrapWithKind(err, errors.KindGraphAPI, "store response to file")
				}

				if err := testCase.Statistics.AddMSErrorResponse(&msResponse, 1); err != nil {
					return errors.Wrap(err, "AddMSErrorResponse")
				}

				filename := fmt.Sprintf("%s_CaseID_%d_repetition_%d_worker_%d_error.json", timestamp, testCase.ID, repetition, workerID)
				if err := storeToFile([]byte(testCase.Log(repetition, workerID, requestMsg, responseMsg)), filename); err != nil {
					return errors.Wrap(err, "store response to file")
				}
				//return errors.New(errors.KindGraphAPI, child.String())
				log.Printf("CaseID_%d_repetition_%d_worker_%d, err: %s", testCase.ID, repetition, workerID, child.String())

			}
		}
	}

	return nil
}

func getGraphApiToken() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "GetConfig")
	}

	//url := "https://login.microsoftonline.com/c3c9dabe-1f88-4098-8b3a-850a742db76a/oauth2/v2.0/token"
	graphApiURL := fmt.Sprintf(`https://login.microsoftonline.com/%s/oauth2/v2.0/token`, cfg.Settings.TenantID)
	body := url.Values{}
	body.Set("client_id", "b372aa3c-9eac-482e-94ff-836c43ea213f")
	body.Set("client_secret", "nFqbtvEd4Uf?x3kSHNn?BaM/[GCdOl95")
	body.Set("grant_type", "client_credentials")
	body.Set("scope", "https://graph.microsoft.com/.default")

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	respBody, err := sendRequest(http.MethodPost, graphApiURL, []byte(body.Encode()), headers)
	if err != nil {
		return errors.Wrap(err, "send request to graph api")
	}

	jsonParsed, err := gabs.ParseJSON(respBody)
	if err != nil {
		return errors.WrapWithKind(err, errors.KindParse, "parse response to json")
	}

	val := jsonParsed.Path("access_token").Data()
	if val == nil {
		return errors.New(errors.KindGraphAPI, jsonParsed.String())
	}
	bearerToken = val.(string)
	return nil
}
