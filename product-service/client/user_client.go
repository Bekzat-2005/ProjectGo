package client

import (
	"fmt"
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
		Get("http://user-service:8080/profile")

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
func GetUserByID(userID uint, token string) (map[string]interface{}, error) {
	var result map[string]interface{}

	resp, err := client.R(). // ✅ дайын client қолдан
					SetHeader("Authorization", "Bearer "+token).
					SetResult(&result).
					Get(fmt.Sprintf("http://user-service:8080/users/%d", userID))

	if err != nil {
		log.Printf("⚠️ Error fetching user: %v", err)
		return nil, err
	}

	log.Printf("📨 Status: %d", resp.StatusCode())
	log.Printf("📨 Body: %s", resp.String())

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("user-service returned %d", resp.StatusCode())
	}

	return result, nil
}
