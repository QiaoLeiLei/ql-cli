package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmd_Execute(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		expected string
	}{
		{
			name:     "help",
			arg:      "--help",
			expected: "一个ql-cli工具集合",
		},
		{
			name:     "version",
			arg:      "--version",
			expected: "ql version 0.0.1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 1. 设置输出缓冲区
			var outBuf, errBuf bytes.Buffer
			rootCmd.SetOut(&outBuf)
			rootCmd.SetErr(&errBuf)

			// 2. 设置参数并执行
			rootCmd.SetArgs([]string{test.arg})
			err := rootCmd.Execute()

			// 3. 验证错误和输出
			if err != nil {
				t.Fatalf("execute() returned unexpected error: %v", err)
			}
			if !strings.Contains(outBuf.String(), test.expected) {
				t.Errorf("output does not contain %q.\nGot: %q", test.expected, outBuf.String())
			}
		})
	}

}
