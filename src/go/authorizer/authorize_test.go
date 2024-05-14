package main

import (
	"fmt"
	"testing"
)

func TestGetWellKnownJwksKeys(t *testing.T) {
	if keys, err := GetWellKnownJwksKeys("us-east-1", "some random text", UserPoolId); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%v", keys)
	}
}