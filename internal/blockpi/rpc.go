package blockpi

import "fmt"

func (c *Client) FirstConfirm() error {
	_, err := c.call("hub_skuFirstConfirm", map[string]interface{}{}, true)
	return err
}

func (c *Client) GenerateApiKey(chainId int, name string) (string, string, error) {
	resp, err := c.call("hub_generateApiKey", map[string]interface{}{
		"chainId": chainId,
		"name":    name,
	}, true)
	if err != nil {
		return "", "", err
	}

	httpRpc, ok1 := resp.Result["httpRpc"].(string)
	wsRpc, ok2 := resp.Result["wsRpc"].(string)

	if !ok1 || !ok2 {
		return "", "", fmt.Errorf("missing rpc fields in result")
	}

	return httpRpc, wsRpc, nil
}
