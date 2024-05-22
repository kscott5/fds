package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
)

var logger *zap.Logger

func getOath2Token(code string) error {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenUrl := os.Getenv("CLIENT_TOKEN_URL")

	clientSecretBasic := fmt.Sprintf("%s:%s", clientId, clientSecret)

	authorization := &strings.Builder{}
	encoder := base64.NewEncoder(base64.StdEncoding, authorization)
	encoder.Write([]byte(clientSecretBasic))
	encoder.Close()

	data := url.Values{}
	data.Add("grant_type","authorization_code")
	data.Add("client_id", clientId)
	data.Add("code", code)
	
	reqBody := strings.NewReader(data.Encode())

	request, err := http.NewRequest(http.MethodPost, tokenUrl, reqBody)
	if err != nil {
		return err
	}

	request.Header.Add("authorization", authorization.String())
	
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	var size int64 = response.ContentLength
	var resBody []byte = make([]byte, size)
	response.Body.Read(resBody)
	return nil
}

func main() {
	logger, _ := zap.NewDevelopment()

	// Cognito Client Pool uses a client secret. A recommendation
	// is keep this data private and safe. Without the use of a
	// private key vault, environment variable are in use and 
	// Lambda funcations are not in use.
	logger.Info("main server")
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
		logger.Debug("localhost processing oauth request")

		httpMethod := strings.ToLower(req.Method)
		query := req.URL.Query()
		
		res.WriteHeader(200)
		res.Header().Add("content-type", "application/json")

		if httpMethod != "get" || query.Has("code") == false {
			res.Write([]byte("invalid server request. view oauth2 and aws cognito client pool redirect url"))
		} else if err := getOath2Token(query.Get("code")); err != nil  {	
			res.Write([]byte(fmt.Sprintf("%s", err)))
		}
	})

	logger.Info("Starting localhost on port: 80")
	logger.Debug("localhost expeccts client id and client secret environment variables")

	http.ListenAndServe(":80", nil)
}