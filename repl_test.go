package main

import "testing"

func TestCleanInput(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "Hello, World!",
			expected: []string{"hello", "world"},
		},
	}

	for _, testCase := range testCases {
		expected := testCase.expected
		actual := cleanInput(testCase.input)

		expectedLen, actualLen := len(actual), len(expected)
		if expectedLen != actualLen {
			t.Errorf("Output length doesn't match expected length:\n\tExpected length of %d\n\tFound length of %d",
				expectedLen,
				actualLen)
		}

		for i := range actual {
			if actual[i] != expected[i] {
				t.Errorf("Output values don't match:\n\tExpected: %s\n\tFound: %s\n", expected[i], actual[i])
			}
		}
	}
}
