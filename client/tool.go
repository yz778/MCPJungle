package client

// ListTools fetches the list of tools, optionally filtered by server name.
//func (c *Client) ListTools(server string) ([]byte, error) {
//	u, _ := url.Parse(c.BaseURL + "/tools")
//	if server != "" {
//		q := u.Query()
//		q.Set("server", server)
//		u.RawQuery = q.Encode()
//	}
//	resp, err := c.HTTPClient.Get(u.String())
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	return io.ReadAll(resp.Body)
//}
