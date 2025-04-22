package tempmail

import (
	"testing"
	"time"
)

func TestNewMailPlusClient(t *testing.T) {
	client, err := NewMailPlusClient()
	if err != nil {
		t.Errorf("NewMailPlusClient() error = %v", err)
		return
	}

	if client.Email == "" {
		t.Error("Expected email to be generated, got empty string")
	}

	if client.client.Timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10s, got %v", client.client.Timeout)
	}
}

func TestMailPlusClient_WaitForConfirmationEmail(t *testing.T) {
	client, _ := NewMailPlusClient()
	client.Email = "fvx44jqru0@any.pink"
	html, err := client.WaitForConfirmationEmail(60 * time.Second)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	t.Log(html)
}

func TestMailPlusClient_FetchMessageHTML(t *testing.T) {
	client, _ := NewMailPlusClient()
	_, err := client.FetchMessageHTML(123)
	if err == nil {
		t.Error("Expected error for non-existent message, got nil")
	}
}
