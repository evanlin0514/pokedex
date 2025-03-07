package main

import(
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string 
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  your mom  ",
			expected: []string{"your", "mom"},
		},
		{
			input:    "  Justin  Biber  ",
			expected: []string{"justin", "biber"},
		},
		{
			input:    "  The Amazing Spiderman  ",
			expected: []string{"the", "amazing", "spiderman"},
		},
		{
			input:    "  GOAT of League  ",
			expected: []string{"goat", "of", "league"},
		},
	}

	passCount := 0
	failCount := 0

	for _, c := range cases {
		test := cleanInput(c.input)
		if len(test) != len(c.expected) {
			fmt.Println(len(test), len(c.expected))
			failCount ++
			t.Errorf(`---------------------------------
			Input slice length: %v
			Expected slice length:	%v
			Fail
			`, len(test), len(c.expected))
			t.FailNow()
		}
		for i, word := range test {
			if word != c.expected[i] {
				failCount ++
				t.Errorf(`---------------------------------
				Expected word:	%v
				Got: %v
				Fail
				`, c.expected[i], word)
				t.FailNow()
			} 
		}
		passCount++
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}