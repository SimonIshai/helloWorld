package main

import (
	e "errors"
	"fmt"
)

func main() {
	fmt.Println("**********************************************")
	fmt.Println()

	err := Deploy()
	if err != nil {
		err = fmt.Errorf("failed to deploy. error: %s", err)
	}
	fmt.Println(err)

	fmt.Printf("\n**********************************************")

}

func Deploy() error {
	err := Test()
	return fmt.Errorf("failed to test %s", err)
}

func Test() error {
	return e.New("failed to publish")
}
