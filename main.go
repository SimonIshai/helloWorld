package main

import (
	"log"
	"time"

	"github.com/SimonIshai/helloWorld/errors"

	"github.com/SimonIshai/helloWorld/config"

	"github.com/SimonIshai/helloWorld/runner"
)

func main() {
	log.Println("starting the Process")

	//testCases := runner.GetTestCases_MaxParallelBatches()
	//testCases := runner.GetTestCases_MaxNumOfRequestsInBatch()
	cfg, err := config.Init("config.yaml")
	//cfg, err := config.GetConfig()
	if err != nil {
		err = errors.Wrap(err, "GetConfig")
		log.Println(err)
		return
	}

	t1 := time.Now()
	if err := runner.Run(cfg.TestCases); err != nil {
		log.Println(err)
		return
	}
	log.Println("Finished processing after", time.Now().Sub(t1).String())
}
