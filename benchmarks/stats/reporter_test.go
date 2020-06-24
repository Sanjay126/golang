package stats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddTagsToName(t *testing.T) {
	tests := []struct {
		name     string
		tags     map[string]string
		expected string
	}{
		{
			name:     "recvd",
			tags:     nil,
			expected: "recvd.no-endpoint",
		},
		{
			name: "recvd",
			tags: map[string]string{
				"endpoint": "hello",
			},
			expected: "recvd.hello",
		},
		{
			name: "r.call",
			tags: map[string]string{
				"host":     "my-host-name",
				"endpoint": "hello",
			},
			expected: "r.call.my-host-name.hello",
		},
		{
			name: "r.call",
			tags: map[string]string{
				"host":     "my\\host.name",
				"endpoint": "hello",
			},
			expected: "r.call.my-host-name.hello",
		},
	}

	for _, testCase := range tests {
		got := addTagsToName(testCase.name, testCase.tags)
		assert.Equal(t, testCase.expected,got)
	}
}