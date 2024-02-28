package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Auth(path, user, pass string) error {
	creds := credentials{user, pass}

	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(creds)
	if err != nil {
		return err
	}

	resp, err := doRequest("POST", path, &body, nil, false)
	if err != nil {
		return err
	}

	return createTokenCache(resp.Body)
}

type cacheToken struct {
	Token string `json:"token"`
}

func createTokenCache(body io.ReadCloser) error {
	token, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	file, err := os.Create(".cacheToken")
	if err != nil {
		return err
	}

	cache := cacheToken{string(token)}
	data, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}

func readCacheToken() (string, error) {
	data, err := os.ReadFile(".cacheToken")
	if err != nil {
		return "", err
	}

	var cache cacheToken

	err = json.Unmarshal(data, &cache)
	if err != nil {
		return "", err
	}

	return cache.Token, nil
}
