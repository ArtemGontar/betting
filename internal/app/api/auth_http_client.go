package api

import (
	"fmt"
	"io"
	"net/http"
)

func HttpRequest(url string, httpVerb string, requestBody io.Reader) []byte {
	client := http.Client{}
	req, err := http.NewRequest(
		httpVerb, url, requestBody,
	)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	return responseBody
}
