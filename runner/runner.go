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

	"github.com/Jeffail/gabs"
	"github.com/SimonIshai/helloWorld/errors"
)

func Run(testCases TestCases) error {
	for i := range testCases {
		err := runCase(&testCases[i])
		if err != nil {
			log.Printf("TestCase %d, err: %s\n", testCases[i].ID, err)
		}
	}
	return nil
}

func runCase(testCase *TestCase) error {

	chanSize := testCase.NumOfRepeats
	ch := make(chan struct{}, chanSize+1)
	for i := 0; i < chanSize; i++ {
		ch <- struct{}{}
	}
	close(ch)
	var wg sync.WaitGroup
	wg.Add(testCase.NumOfConcurrentWorkers)
	for i := 0; i < testCase.NumOfConcurrentWorkers; i++ {
		go func() {
			defer wg.Done()

			for range ch {
				resp, err := sendRequest(testCase)
				if err != nil {
					log.Println("runCase worker err:", errors.Wrap(err, errors.KindHttp, "send request"))
					return
				}

				err = processResponse(testCase, resp)
				if err != nil {
					log.Println("runCase worker err:", errors.Wrap(err, errors.KindHttp, "process response"))
					return
				}
			}

		}()
	}
	wg.Wait()
	return nil
}

func sendRequest(testCase *TestCase) ([]byte, error) {

	var url string
	var method string
	var err error
	var bodyReader *bytes.Reader
	var message string
	if testCase.IsBatch {
		method = testCase.BatchRequest.Method
		url = urlAPI + testCase.BatchRequest.URL
		body, err := json.MarshalIndent(testCase.BatchRequest.Body, "", "\t")
		if err != nil {
			return nil, errors.Wrap(err, errors.KindParse, "json marshal")
		}
		message += method + ": " + url + "\n" + string(body) + "\n"
		bodyReader = bytes.NewReader(body)

	} else {
		url = urlAPI + testCase.Request.URL
		method = testCase.Request.Method
		if testCase.Request.Body != "" {
			bodyReader = bytes.NewReader([]byte(testCase.Request.Body))
		}
		message += method + ": " + url + "\n" + testCase.Request.Body + "\n"
	}
	testCase.RequestMsg = message

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, errors.KindHttp, "NewRequest")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errors.KindHttp, "client.Do")
	}
	defer resp.Body.Close() //nolint: errcheck

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, errors.KindHttp, "read response body")
	}

	return respBody, nil
}

func processResponse(testCase *TestCase, resp []byte) error {

	testCase.ResponseMsg = string(resp)

	//filename := time.Now().Format("2006-01-02T15:04") + ".json"
	//if err := storeToFile(resp, filename); err != nil {
	//	return errors.Wrap(err, errors.KindFileSystem, "store response to file")
	//}

	container, err := gabs.ParseJSON(resp)
	if err != nil {
		return errors.Wrap(err, errors.KindParse, "parse response to json")
	}

	statusOKPath := "status"
	if elem := container.Path(statusOKPath); elem != nil {
		httpStatus := elem.Data().(int)
		if httpStatus != http.StatusOK {

		}
	}

	errPath := "error"
	//errCodePath := "error.code"
	//errMessagePath := "error.message"

	if elem := container.Path(errPath); elem != nil {

		timestamp := time.Now().Format("2006-01-02T15:04")
		filename := fmt.Sprintf("%s_CaseID_%d_error.json", timestamp, testCase.ID)

		if err := storeToFile([]byte(testCase.String()), filename); err != nil {
			return errors.Wrap(err, errors.KindFileSystem, "store response to file")
		}
	}

	return nil
}

func storeToFile(bytesToStore []byte, filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, errors.KindFileSystem, "open file")
	}

	defer f.Close()
	bytesToStore = append(bytesToStore, []byte("\n")...)
	if _, err = f.Write(bytesToStore); err != nil {
		return errors.Wrap(err, errors.KindFileSystem, "file write")
	}
	return nil
}
