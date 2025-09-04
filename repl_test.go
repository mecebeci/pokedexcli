package main

import (
	"testing"
)

func TestClenInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: " mec mec mec",
			expected: []string{"mec", "mec","mec"},
		},
	}

	for _,c := range cases{
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected){
			t.Errorf("actual and expected length doesn't match")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord{
				t.Errorf("word and expected word doesn't match")
			}
		}
	}
}