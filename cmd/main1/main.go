package main

import (
	"fmt"
	e "github.com/SimonIshai/helloWorld/errors"
)

func main() {
	fmt.Println("**********************************************")
	fmt.Println()

	err := Deploy()
	if err != nil {
		err = e.Wrap(err, 0, "deploying")
	}
	fmt.Println(err)

	fmt.Printf("\n**********************************************")

}

func Deploy() error {
	err := Test()
	return e.Wrap(err, 0, "testing")
}

func Test() error {
	return e.New(e.KindKafka, "failed to publish")
}
