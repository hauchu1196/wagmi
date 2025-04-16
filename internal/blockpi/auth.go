package blockpi

import (
	"fmt"
	"time"

	"github.com/hauchu1196/wagmi/internal/proxy"
)

func (c *Client) Register(email, password string) error {
	var lastErr error

	// Get list of working proxies once
	fmt.Println("🔍 Đang tìm proxy hoạt động...")
	workingProxies, err := proxy.GetLiveProxy(5 * time.Second)
	if err != nil {
		return fmt.Errorf("❌ Không tìm được proxy hoạt động: %v", err)
	}
	fmt.Printf("✅ Tìm thấy %d proxy hoạt động\n", len(workingProxies))

	// Try each working proxy until one succeeds
	for _, proxyURL := range workingProxies {
		fmt.Println("✅ Thử proxy:", proxyURL)
		proxyClient := proxy.MakeClient(proxyURL, 10*time.Second)
		_, err := c.callWithClient("hub_register", map[string]interface{}{
			"email":           email,
			"password":        password,
			"confirmPassword": password,
		}, proxyClient, false)

		if err == nil {
			return nil
		}

		lastErr = err
		fmt.Printf("❌ Đăng ký thất bại với proxy %s: %v\n", proxyURL, err)
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("❌ Đăng ký thất bại với tất cả proxy: %v", lastErr)
}

func (c *Client) EmailConfirm(code string) (string, error) {
	resp, err := c.call("hub_emailConfirm", map[string]interface{}{
		"code": code,
	}, false)
	if err != nil {
		return "", err
	}
	return resp.Result["token"].(string), nil
}
