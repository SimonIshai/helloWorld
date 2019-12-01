package runner

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/SimonIshai/helloWorld/config"

	"github.com/SimonIshai/helloWorld/errors"
)

func sendRequest(method, url string, body []byte, headers map[string]string) ([]byte, error) {

	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, errors.WrapWithKind(err, errors.KindHttp, "NewRequest")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

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

	//logFolder := "test-case-logs"
	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "GetConfig")
	}

	logFolder := cfg.Settings.LogFolder
	_, err = os.Stat(logFolder)
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
