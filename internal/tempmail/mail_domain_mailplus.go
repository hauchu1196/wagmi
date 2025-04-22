package tempmail

import (
	"fmt"
	"time"
)

type MailDomainMailPlusClient struct {
	MailPlusClient *MailPlusClient
	OriginalEmail  string
}

func NewMailDomainMailPlusClient() (*MailDomainMailPlusClient, error) {
	domain := "hauchu.space"
	targetEmail := "geaba@mailto.plus"

	randomID := genRandomString(10)
	email := fmt.Sprintf("%s@%s", randomID, domain)
	mailPlusClient, err := NewMailPlusClient()
	if err != nil {
		return nil, err
	}
	mailPlusClient.Email = targetEmail
	return &MailDomainMailPlusClient{
		MailPlusClient: mailPlusClient,
		OriginalEmail:  email,
	}, nil
}

func (c *MailDomainMailPlusClient) WaitForConfirmationEmail(timeout time.Duration) (string, error) {
	return c.MailPlusClient.WaitForConfirmationEmail(timeout)
}

func (c *MailDomainMailPlusClient) GetEmail() string {
	return c.OriginalEmail
}

func (c *MailDomainMailPlusClient) GetPassword() string {
	return ""
}
