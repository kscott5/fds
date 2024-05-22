package main

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	logger.Info("main server")
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
		logger.Debug("localhost processing oauth request")

		res.Write([]byte("<h1>Keep it simple</h1>"))
		res.Write([]byte("The address bar contains either the access_token and id_token or error message<br/><br/>"))
		res.Write([]byte("Access API gateway with curl and access_token or id_token.<br/><br/>"))
		res.Write([]byte("NOTE: the hash tag is reserve use with web or SPA apps. <br/>"))
		res.Write([]byte("Example: javascript:document.location.hash"))
	})

	logger.Info("Starting localhost on port: 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
}
