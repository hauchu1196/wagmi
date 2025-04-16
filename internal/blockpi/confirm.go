package blockpi

import (
	"fmt"
	"strings"
)

func ExtractConfirmationCode(emailBody string) (string, error) {
	start := strings.Index(emailBody, "activate-account?code=")
	if start == -1 {
		return "", fmt.Errorf("Không tìm thấy mã xác nhận trong email")
	}
	start += len("activate-account?code=")

	end := start
	for end < len(emailBody) && isHexChar(emailBody[end]) {
		end++
	}

	return emailBody[start:end], nil
}

func isHexChar(c byte) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'f') ||
		(c >= 'A' && c <= 'F') ||
		c == 'x'
}
