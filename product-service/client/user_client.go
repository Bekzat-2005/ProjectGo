package client

import (
	"github.com/go-resty/resty/v2"
	"log"
)

var client = resty.New()

func init() {
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.Printf("[Resty] Request: %s %s", req.Method, req.URL)
		return nil
	})
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log.Printf("[Resty] Response: %d %s", resp.StatusCode(), resp.Request.URL)
		return nil
	})
}

// GetUserProfile делает запрос к user-service по /profile
func GetUserProfile(token string) (string, error) {
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get("http://localhost:8080/profile")

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
