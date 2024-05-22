package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
)

func getOath2Token(code string) ([]byte, error) {
	logger, _ := zap.NewDevelopment()
	logger.Info(fmt.Sprintf("get oauth2 token with code: %s", code))

	if code == "" {
		return nil, fmt.Errorf("requires authorization code received from the authorization server")
	}

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenUrl := os.Getenv("CLIENT_TOKEN_URL")
	redirectUri := os.Getenv("CLIENT_REDIRECT_URI")
	
	clientSecretBasic := fmt.Sprintf("%s:%s", clientId, clientSecret)
	
	authorization := &strings.Builder{}
	encoder := base64.NewEncoder(base64.StdEncoding, authorization)
	
	if _, err := encoder.Write([]byte(clientSecretBasic)); err != nil {
		return nil, err
	}
	encoder.Close()

	data := url.Values{}
	data.Add("grant_type","authorization_code")
	data.Add("client_id", clientId)
	data.Add("code", code)
	data.Add("redirect_uri", redirectUri)
	
	body := strings.NewReader(data.Encode())
	
	logger.Debug(fmt.Sprintf("creating http new request body: %s", data.Encode()))
	request, err := http.NewRequest(http.MethodPost, tokenUrl, body)
	if err != nil {
		return nil, err
	}

	request.URL.Scheme = "http"
	request.Header.Add("content-type", "application/www-form-urlencoded")
	request.Header.Add("authorization", fmt.Sprintf("Basic %s", authorization))
	
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	logger.Debug("athorization access token request completed")
	var size int64 = response.ContentLength
	var tokenData []byte = make([]byte, size)
	if _, err := response.Body.Read(tokenData); err != nil {
		return nil, err
	}
	
	logger.Debug(fmt.Sprintf("token data %s", string(tokenData)))
	return tokenData, nil
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

		if httpMethod != "get" || !query.Has("code") {
			res.Write([]byte(`{"invalid_request": "requires authorization code received from the authorization server"}`))
		} 
		
		code := query.Get("code")
		if token, err := getOath2Token(code); err != nil { 
			res.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err)))
		} else {
			res.Write(token)
		}
	})

	logger.Info("Starting localhost on port: 8080")
	logger.Debug("localhost expects client id and client secret environment variables")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
}