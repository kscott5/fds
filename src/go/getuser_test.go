package main

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestGetUser(t *testing.T) {
	logger, _ = zap.NewDevelopment()
	pathParams := make(map[string]interface{}, 1)
	pathParams["_id"] = "97ed8408f0e24422ba619884ab7d116d"

	request := Request{
		HttpMethod: "GET",
		Resource: "/users/{_id}",
		PathParameters: pathParams,
	}

	if response, err := getUser(context.Background(), &request); err != nil {
		t.Error(err)
	} else {
		fmt.Println(response)
	}
}