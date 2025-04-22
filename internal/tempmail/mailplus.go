package tempmail

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// --- MailTMClient manages temporary email lifecycle ---

type MailPlusClient struct {
	Email    string
	Password string
	Token    string
	client   *http.Client
}

const mailPlusBaseURL = "https://tempmail.plus/api"

var MailPlusDomains = []string{"mailto.plus", "fexpost.com", "fexbox.org", "mailbox.in.ua", "rover.info", "chitthi.in", "fextemp.com", "any.pink", "merepost.com"}

// --- Public: Create client with random email & fetch JWT token ---

func NewMailPlusClient() (*MailPlusClient, error) {
	randomID := genRandomString(10)

	// Get random domain
	domain := MailPlusDomains[rand.Intn(len(MailPlusDomains))]
	email := fmt.Sprintf("%s@%s", randomID, domain)

	return &MailPlusClient{
		Email:  email,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// --- Wait until a message arrives and return HTML body ---

func (c *MailPlusClient) WaitForConfirmationEmail(timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		req, _ := http.NewRequest("GET", mailPlusBaseURL+"/mails", nil)
		q := req.URL.Query()
		q.Add("email", c.Email)
		q.Add("epin", c.Password)
		q.Add("limit", "20")
		req.URL.RawQuery = q.Encode()

		resp, err := c.client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var messages struct {
			MailList []struct {
				IsNew  bool `json:"is_new"`
				MailId int  `json:"mail_id"`
			} `json:"mail_list"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
			return "", err
		}

		for _, message := range messages.MailList {
			if message.IsNew {
				// Read the latest message
				msgID := message.MailId
				return c.FetchMessageHTML(msgID)
			}
		}
		
		time.Sleep(2 * time.Second)
	}

	return "", fmt.Errorf("timeout: không nhận được email trong %s", timeout)
}

// --- Fetch HTML content of an email by ID ---

func (c *MailPlusClient) FetchMessageHTML(msgID int) (string, error) {
	url := fmt.Sprintf("%s/mails/%d", mailPlusBaseURL, msgID)
	req, _ := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("email", c.Email)
	q.Add("epin", c.Password)
	req.URL.RawQuery = q.Encode()

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

func (c *MailPlusClient) GetEmail() string {
	return c.Email
}

func (c *MailPlusClient) GetPassword() string {
	return c.Password
}