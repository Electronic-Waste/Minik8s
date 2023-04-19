package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNewCmdVersion(t *testing.T) {
	// test the build of a new command by run the command to see if the error is caused
	var buf bytes.Buffer
	cmd := NewCmdVersion(&buf)
	err := cmd.Execute()
	if err != nil {
		t.Error("can't make a new version cmd\n")
	}
}

func TestRunVersion(t *testing.T) {
	var buf bytes.Buffer
	iface := make(map[string]interface{})
	flagNameOutput := "output"
	cmd := NewCmdVersion(&buf)

	testCases := []struct {
		name              string
		flag              string
		expectedError     bool
		shouldBeValidYAML bool
		shouldBeValidJSON bool
	}{
		{
			name: "valid: run without flags",
		},
		{
			name:              "valid: run with flag 'yaml'",
			flag:              "yaml",
			shouldBeValidYAML: true,
		},
		{
			name:              "valid: run with flag 'json'",
			flag:              "json",
			shouldBeValidJSON: true,
		},
		{
			name:          "invalid: run with unsupported flag",
			flag:          "unsupported-flag",
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			if len(tc.flag) > 0 {
				if err = cmd.Flags().Set(flagNameOutput, tc.flag); err != nil {
					goto error
				}
			}
			buf.Reset()
			if err = RunVersion(&buf, cmd); err != nil {
				goto error
			}
			if buf.String() == "" {
				err = errors.New("empty output")
				goto error
			}
			if tc.shouldBeValidYAML {
				if !strings.Contains(buf.String(), "yaml") {
					err = errors.New("can't parse yaml cmd")
				}
			} else if tc.shouldBeValidJSON {
				if !strings.Contains(buf.String(), "json") {
					err = errors.New("can't parse json cmd")
				}
			}
		error:
			if (err != nil) != tc.expectedError {
				t.Errorf("Test case %q: RunVersion expected error: %v, saw: %v; %v", tc.name, tc.expectedError, err != nil, err)
			}
		})
	}
}
