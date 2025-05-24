package client

import "io"

// ListServers fetches the list of registered servers.
func (c *Client) ListServers() ([]byte, error) {
	u, _ := c.constructAPIEndpoint("/servers")
	resp, err := c.HTTPClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// DeregisterServer deletes a server by name.
//func (c *Client) DeregisterServer(name string) error {
//	req, _ := http.NewRequest(http.MethodDelete, c.BaseURL+"/servers/"+name, nil)
//	resp, err := c.HTTPClient.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusNoContent {
//		body, _ := io.ReadAll(resp.Body)
//		return fmt.Errorf("unexpected status %s, body: %s", resp.Status, body)
//	}
//	return nil
//}
//
//// InvokeTool sends a JSON payload to invoke a tool.
//func (c *Client) InvokeTool(payload map[string]any) ([]byte, error) {
//	body, _ := json.Marshal(payload)
//	resp, err := c.HTTPClient.Post(c.BaseURL+"/tools/invoke", "application/json", bytes.NewReader(body))
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	return io.ReadAll(resp.Body)
//}
