package forismatic

import (
	"encoding/json"
	"io"
	"motivation-bot/logging"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	logger *logging.Logger
}

const BaseUrl = "https://api.forismatic.com/api/1.0/"

type GetQuoteResponse struct {
	QuoteText   string `json:"quoteText"`
	QuoteAuthor string `json:"quoteAuthor"`
	SenderName  string `json:"senderName"`
	SenderLink  string `json:"senderLink"`
	QuoteLink   string `json:"quoteLink"`
}

func NewClient(logger *logging.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

func (c *Client) GetQuote(lang string, key ...int) GetQuoteResponse {
	formData := url.Values{
		"method": []string{"getQuote"},
		"format": []string{"json"},
		"lang":   []string{lang},
	}

	if len(key) > 0 {
		keyStr := strconv.Itoa(key[0])
		formData.Add("key", keyStr)
	}

	resp, err := http.PostForm(BaseUrl, formData)
	if err != nil {
		c.logger.Errorf("Error during request: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Errorf("Error on parsing response: %s", err)
	}

	var quoteResponse GetQuoteResponse
	err = json.Unmarshal(body, &quoteResponse)
	if err != nil {
		c.logger.Errorf("Error on Unmarshal response: %s", err)
	}
	return quoteResponse
}
