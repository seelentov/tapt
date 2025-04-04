package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type getTextRequest struct {
	Content string `json:"content"`
}

func GetText() (string, error) {
	resp, err := http.Get("http://api.quotable.io/random")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(bb))
	}

	var req getTextRequest

	if err := json.Unmarshal(bb, &req); err != nil {
		return "", nil
	}

	return req.Content, nil
}
