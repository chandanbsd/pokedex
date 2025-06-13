package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},

		{
			input:    "  Hello  World!  ",
			expected: []string{"hello", "world!"},
		},

		{
			input:    "  Mad ad 2 % !!!  World!  ",
			expected: []string{"mad", "ad", "2", "%", "!!!", "world!"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("Length of output does not match input")
			return
		}

		for i := range actual {

			if actual[i] != c.expected[i] {
				t.Errorf("The actual word does not match expected word")
				return
			}
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
		}
	}
}
