package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders_Parse(t *testing.T) {
	testTable := []struct {
		name         string
		headersLine  string
		expected     Headers
		expectedDone bool
		expectedN    int
		expectedErr  bool
	}{
		{
			name:         "single header",
			headersLine:  "Host: localhost:42069\r\n",
			expected:     Headers{"host": "localhost:42069"},
			expectedDone: false,
			expectedN:    len("Host: localhost:42069\r\n"),
		},
		//{
		//	name:         "multiple header keys",
		//	headersLine:  "Host: localhost:42069\r\nHost: localhost:42068\r\n\r\n",
		//	expected:     Headers{"host": "localhost:42069, localhost:42068"},
		//	expectedDone: false,
		//	expectedN:    len("Host: localhost:42069\r\nHost: localhost:42068\r\n\r\n"),
		//},
		{
			name:        "whitespaces before key",
			headersLine: "     Host: localhost:42069\r\n",
			expectedErr: true,
		},
		{
			name:        "invalid char in header key",
			headersLine: "H©st: localhost:42069\r\n",
			expectedErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			headers := New()
			n, done, err := headers.Parse([]byte(testCase.headersLine))
			if testCase.expectedErr {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, headers)
			assert.Equal(t, testCase.expectedN, n)
			assert.Equal(t, testCase.expectedDone, done)
		})
	}
}
