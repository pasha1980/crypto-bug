package main

import (
	"crypto-bug/config"
	"crypto-bug/test"
	"fmt"
)

type TestInterface interface {
	Do() bool
}

var tests = []TestInterface{
	test.CacheTest{},
}

const finishMessage = `
Testing completed. 
Successfully completed %d out of %d tests
`

func main() {
	var success int64
	config.Initialization()
	fmt.Println("Testing...")
	for _, test := range tests {
		if test.Do() {
			success++
		}
	}
	message := fmt.Sprintf(finishMessage, success, len(tests))
	fmt.Println(message)
}
