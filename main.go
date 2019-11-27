package main

import (
	"log"

	"github.com/SimonIshai/helloWorld/runner"
)

func main() {
	log.Println("starting the Process")

	testCases := runner.GetTestCases()

	if err := runner.Run(testCases); err != nil {
		log.Println(err)
	}
}
