package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	logger.Info("main server")
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
		logger.Debug("localhost processing oauth request")

		res.WriteHeader(200)
		res.Header().Add("content-type", "application/json")

		body := make(map[string]interface{},1)
		body["status"] = 200

		query := req.URL.Query()
		if query.Has("error") && query.Has("error_description") {
			body["error"] = query.Get("error")
			body["error_description"] = query.Get("error_description")

			buffer, _ := json.Marshal(body)
			res.Write(buffer)
			return
		}
		
		data := make(map[string]string,1)
		if query.Has("id_token")  {			
			data["id_token"] = query.Get("id_token")
		}

		if query.Has("access_token") {
			data["access_token"] = query.Get("access_token")
		}

		if query.Has("token_type") {
			data["token_type"] = query.Get("token_type")
		}

		if query.Has("expires_in") {
			data["expires_in"] = fmt.Sprintf("%s sec", query.Get("expires_in"))
		}

		if query.Has("state") {
			data["state"] = query.Get("state")
		}

		if len(data) == 5 {
			body["message"] = "invalid response from oauth2/authorize request"
		} else { 
			body["message"] = "success"
			body["data"] = data
		}

		buffer, _ := json.Marshal(body)
		res.Write(buffer)
	})

	logger.Info("Starting localhost on port: 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
}