package cmd

import (
	`context`
	`net/http`
)

func request(ctx *context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(*ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}
