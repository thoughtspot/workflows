package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func Test_Prompt(t *testing.T) {
	tests := []struct {
		input    string
		question string
		message  string
		want     string
	}{
		{
			input:    "github-auth-token\n",
			question: "Github Auth Token",
			message:  "",
			want:     "github-auth-token",
		},
	}

	for _, test := range tests {
		reader := strings.NewReader(test.input)
		var writer bytes.Buffer

		rw := NewRW(reader, &writer)
		got := rw.Prompt(test.question, test.message)

		if got != test.want {
			t.Errorf("\nGot: %v\nWant: %v\n", got, test.want)
		}

	}
}
