package model

import (
	"fmt"
	"sync"
	"time"

	"github.com/SimonIshai/helloWorld/errors"
)

type Statistics struct {
	SpanStartTime      time.Time    `json:"-"`
	SpanEndTime        time.Time    `json:"-"`
	TimeSpanSec        float32      `json:"time_span_sec"`
	MSErrorResponses   []MSResponse `json:"-"`
	TotalRequests      int          `json:"total_requests"`
	TotalRequestErrors int          `json:"total_request_errors"`
	ErrorPercentage    float32      `json:"errors_percentage"`
}

type MSResponse struct {
	ReqID      string `json:"id"`
	HttpStatus int    `json:"status"`
	Body       struct {
		MSError MSError `json:"error"`
	} `json:"body"`
}

type MSError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var mutex1 = &sync.Mutex{}
var mutex2 = &sync.Mutex{}

func (s *Statistics) CalcErrorPercentage() error {

	mutex1.Lock()
	defer mutex1.Unlock()

	if s.TotalRequests <= 0 {
		return errors.New(errors.KindSystem, "TotalRequests is 0")
	}
	if s.TotalRequestErrors == 0 {
		s.ErrorPercentage = 0
		return nil
	}
	s.ErrorPercentage = (float32(s.TotalRequestErrors) / float32(s.TotalRequests)) * 100
	return nil
}
func (s *Statistics) CalcTimeSpan() error {

	s.TimeSpanSec = float32(s.SpanEndTime.Sub(s.SpanStartTime).Milliseconds()) / 1000
	return nil
}

func (s *Statistics) AddMSErrorResponse(msResponse *MSResponse, numOfRequestErrors int) error {
	mutex2.Lock()
	defer mutex2.Unlock()

	s.MSErrorResponses = append(s.MSErrorResponses, *msResponse)
	s.TotalRequestErrors += numOfRequestErrors
	err := s.CalcErrorPercentage()
	if err != nil {
		return errors.Wrap(err, "CalcErrorPercentage")
	}
	return nil
}

func (s *Statistics) Refresh() {
	s.CalcErrorPercentage()
	s.CalcTimeSpan()
}

func (s *Statistics) String() string {

	s.Refresh()

	return fmt.Sprintf(`
Total Requests %d
Error Responses %d
Errors Rate %.2f %%
Time span %v
`, s.TotalRequests, len(s.MSErrorResponses), s.ErrorPercentage, s.SpanEndTime.Sub(s.SpanStartTime))
}
