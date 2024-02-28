package requests

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func doRequest(method, path string, body io.Reader, headers map[string]string, auth bool) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:8090%s", path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if _, ok := headers["Content-Type"]; !ok {
		req.Header.Add("Content-Type", "application/json")
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	if auth {
		token, err := readCacheToken()
		if err != nil {
			log.Println("Não foi possível obter o token")
			return nil, err
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	return http.DefaultClient.Do(req)
}
