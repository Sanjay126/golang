package internal_test

import (
	"bitbucket.org/dreamplug-backend/benchmarks/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsGopher(t *testing.T) {
	tests:=[]struct{
		email string
		expectedName string
		isGopher bool
	}{
		{
			email:"sanjay@gmail.com",
			expectedName:"",
			isGopher:false,
		},
		{
			email:"sanjay@golang.org",
			expectedName:"sanjay",
			isGopher:true,
		},
	}
	for _,testCase:=range tests{
		name,gopher:=internal.IsGopher(testCase.email)
		assert.Equal(t,testCase.expectedName,name)
		assert.Equal(t,testCase.isGopher,gopher)
	}
}

func BenchmarkIsGopher(b *testing.B) {
	for i := 0; i < b.N; i++ {
		internal.IsGopher("sanjay@golang.org")

	}
}