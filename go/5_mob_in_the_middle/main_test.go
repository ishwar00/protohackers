package main

import (
	"fmt"
	"testing"
)


func Test_bogusCoinAddress(t *testing.T) {
	tests := []struct {
		input string
		expOutput string
	}  {
		{
			input: "Please send the payment of 750 Boguscoins to 7Y7wBRhrKrQ9s9GX4L2EzYK0jmqlC",
			expOutput: "Please send the payment of 750 Boguscoins to 7YWHMfk9JZe0LM0g1ZauHuiSxhI",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			output := WriteBoguscoinAddress(tt.input)
			if output != tt.expOutput {
				t.Errorf("output: %s, expOutput: %s", output, tt.expOutput)
			}
		})
	}
}
