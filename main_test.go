package main

import (
	"fmt"
	"testing"
)

var testCases = []struct {
	in  string
	out string
}{
	{"", "a"},
	{"a8", "a9"},
	{"aii", "aij"},
	{"asdsa9", "asdsba"},
}

func TestCodeGeneration(t *testing.T) {

	for _, testCase := range testCases {
		expected := testCase.out
		response, _ := GenerateNextCode(testCase.in)
		if response != expected {
			fmt.Println("Expected :" + expected + " Received :" + response)
			t.Error("Test failed")
		}
	}
}
