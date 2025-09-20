package main

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	type testCase struct {
		args        []string
		expectedCfg Config
		expectedErr string
	}

	runCases := []testCase{
		{
			args:        []string{"-input", "orders.csv"},
			expectedCfg: Config{Input: "orders.csv", Limit: 100, Verbose: false, Format: "text"},
			expectedErr: "",
		},
		{
			args:        []string{"-input", "log.txt", "-limit", "10", "-verbose", "-format", "json"},
			expectedCfg: Config{Input: "log.txt", Limit: 10, Verbose: true, Format: "json"},
			expectedErr: "",
		},
	}

	submitCases := append(runCases, []testCase{
		{
			args:        []string{},
			expectedCfg: Config{},
			expectedErr: "missing -input",
		},
		{
			args:        []string{"-input", "data.txt", "-limit", "0"},
			expectedCfg: Config{},
			expectedErr: "invalid -limit: must be between 1 and 1000",
		},
		{
			args:        []string{"-input", "data.txt", "-format", "xml"},
			expectedCfg: Config{},
			expectedErr: "invalid -format: must be 'text' or 'json'",
		},
		{
			args:        []string{"-input", "bulk.csv", "-limit", "1000", "-format", "text"},
			expectedCfg: Config{Input: "bulk.csv", Limit: 1000, Verbose: false, Format: "text"},
			expectedErr: "",
		},
	}...)

	testCases := runCases
	if withSubmit {
		testCases = submitCases
	}

	skipped := len(submitCases) - len(testCases)

	passCount := 0
	failCount := 0

	for _, tc := range testCases {
		cfg, err := parseArgs(tc.args)
		if tc.expectedErr == "" {
			if err != nil || cfg != tc.expectedCfg {
				failCount++
				errMsg := "<nil>"
				if err != nil {
					errMsg = err.Error()
				}
				t.Errorf(`---------------------------------
Input args: %v

Expected Config: %+v
Expected Error:  <nil>

Actual Config:   %+v
Actual Error:    %s
Fail
`, tc.args, tc.expectedCfg, cfg, errMsg)
			} else {
				passCount++
				fmt.Printf(`---------------------------------
Input args: %v

Expected Config: %+v
Actual Config:   %+v
Pass
`, tc.args, tc.expectedCfg, cfg)
			}
		} else {
			if err == nil || err.Error() != tc.expectedErr {
				failCount++
				errMsg := "<nil>"
				if err != nil {
					errMsg = err.Error()
				}
				t.Errorf(`---------------------------------
Input args: %v

Expected Error: %s
Actual Error:   %s
Fail
`, tc.args, tc.expectedErr, errMsg)
			} else {
				passCount++
				fmt.Printf(`---------------------------------
Input args: %v

Expected Error: %s
Actual Error:   %s
Pass
`, tc.args, tc.expectedErr, err.Error())
			}
		}
	}

	fmt.Println("---------------------------------")
	if skipped > 0 {
		fmt.Printf("%d passed, %d failed, %d skipped\n", passCount, failCount, skipped)
	} else {
		fmt.Printf("%d passed, %d failed\n", passCount, failCount)
	}
}

// withSubmit is set at compile time depending
// on which button is used to run the tests
var withSubmit = true
