package telegram

import "net/http"

type Client struct {
	host string
	basePatch string
	client http.Client
}

func New(host string, token string) {
	return Client{
		host: host,
		basePatch: newBasePath(token),
		client: http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(ofset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	// do reuest <- getUpdates
}
func(c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u:=url.URL{
		Scheme: "https",
		Host: c.host,
		Path: path.Join(c.basePatch, method),
	}
	req,err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	
	req.URL.RawQuery = query.Encode()
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	
}

func (c *Client) SendMessage() {

}