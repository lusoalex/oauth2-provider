package oauth2Provider

import (
	"strings"
	"testing"
)

type DataTestCase struct {
	name   string
	method string
	url    string
	status int
	error  string
}

func assertStatus(t *testing.T, testCase *DataTestCase, status int) {
	if status != testCase.status {
		t.Errorf("%v test : Got unexpected status code (%v) while expecting (%v)",
			testCase.name, status, testCase.status)
	}
}

func assertBodyErrorMessage(t *testing.T, testCase *DataTestCase, body string) {
	if !strings.Contains(body, testCase.error) {
		t.Errorf("%v test : Body response does not contains expected error message ([%v]), got ([%v])",
			testCase.name, testCase.error, body)
	}
}
