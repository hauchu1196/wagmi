package tempmail

import "time"

// MailDomain represents a temporary email service provider interface
type TempMail interface {
	// WaitForConfirmationEmail waits for a confirmation email to arrive within the given timeout duration
	// Returns the HTML content of the email or an error if timeout is reached
	WaitForConfirmationEmail(timeout time.Duration) (string, error)

	// GetEmail returns the generated email address
	GetEmail() string

	// GetPassword returns the password for the email account, if applicable
	GetPassword() string
}
