package tempmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// --- Structs for mail.tm API ---

type Account struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Message struct {
	ID   string `json:"id"`
	From struct {
		Address string `json:"address"`
		Name    string `json:"name"`
	} `json:"from"`
	Subject string `json:"subject"`
}

type MessageContent struct {
	Text string          `json:"text"`
	HTML json.RawMessage `json:"html"`
}

// --- MailTMClient manages temporary email lifecycle ---

type MailTMClient struct {
	Email    string
	Password string
	Token    string
	client   *http.Client
}

const baseURL = "https://api.mail.tm"

// --- Public: Create client with random email & fetch JWT token ---

func NewClient() (*MailTMClient, error) {
	randomID := genRandomString(10)
	password := genRandomString(12)

	// Get random domain
	domain, err := fetchRandomDomain()
	if err != nil {
		return nil, err
	}
	email := fmt.Sprintf("%s@%s", randomID, domain)

	// Register new account
	acc := Account{Address: email, Password: password}
	body, _ := json.Marshal(acc)
	_, err = http.Post(baseURL+"/accounts", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("tạo tài khoản thất bại: %w", err)
	}

	// Login to get JWT
	loginBody, _ := json.Marshal(map[string]string{
		"address":  email,
		"password": password,
	})
	resp, err := http.Post(baseURL+"/token", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		return nil, fmt.Errorf("login thất bại: %w", err)
	}
	defer resp.Body.Close()

	var tokenRes TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}

	return &MailTMClient{
		Email:    email,
		Password: password,
		Token:    tokenRes.Token,
		client:   &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// --- Wait until a message arrives and return HTML body ---

func (c *MailTMClient) WaitForConfirmationEmail(timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		req, _ := http.NewRequest("GET", baseURL+"/messages", nil)
		req.Header.Set("Authorization", "Bearer "+c.Token)

		resp, err := c.client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var messages struct {
			HydraMember []Message `json:"hydra:member"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
			return "", err
		}

		if len(messages.HydraMember) > 0 {
			// Read the latest message
			msgID := messages.HydraMember[0].ID
			return c.FetchMessageHTML(msgID)
		}

		time.Sleep(2 * time.Second)
	}

	return "", fmt.Errorf("timeout: không nhận được email trong %s", timeout)
}

// --- Fetch HTML content of an email by ID ---

func (c *MailTMClient) FetchMessageHTML(msgID string) (string, error) {
	url := fmt.Sprintf("%s/messages/%s", baseURL, msgID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var content MessageContent
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return "", err
	}

	// Try to unmarshal as string first
	var htmlStr string
	if err := json.Unmarshal(content.HTML, &htmlStr); err == nil {
		return htmlStr, nil
	}

	// If not a string, try to unmarshal as array and take first element
	var htmlArr []string
	if err := json.Unmarshal(content.HTML, &htmlArr); err == nil && len(htmlArr) > 0 {
		return htmlArr[0], nil
	}

	return "", fmt.Errorf("không thể parse nội dung email")
}

func (c *MailTMClient) GetEmail() string {
	return c.Email
}

func (c *MailTMClient) GetPassword() string {
	return c.Password
}

// --- Internal: Random string + domain fetch ---

func fetchRandomDomain() (string, error) {
	resp, err := http.Get(baseURL + "/domains")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		HydraMember []struct {
			Domain string `json:"domain"`
		} `json:"hydra:member"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.HydraMember) == 0 {
		return "", fmt.Errorf("không có domain khả dụng")
	}

	return res.HydraMember[0].Domain, nil
}
