package main

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestGetWellKnownJwksKeys(t *testing.T) {
	logger, _ = zap.NewDevelopment()
	region := os.Getenv("AWS_REGION")

	if keys, err := GetWellKnownJwksKeys(region, UserPoolId); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%v", keys)
	}
}

func TestValidateAuthToken(t *testing.T) {
	logger, _ = zap.NewDevelopment()

	region := os.Getenv("AWS_REGION")
	authToken := os.Getenv("AWS_SESSION_TOKEN")

	if token, err := ValidateAuthToken(region, authToken); err != nil {
		t.Error(err)
 	} else {
		t.Log(token)
	}


}