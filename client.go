package rusender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "https://api.beta.rusender.ru/api/v1/external-mails"
const sendByTemplateUrl = "/send-by-template"
const sendUrl = "/send"

type Client struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

func NewClient(client *http.Client, apiKey string) *Client {
	return &Client{
		client:  client,
		apiKey:  apiKey,
		baseURL: baseUrl,
	}
}

func (c *Client) Send(ctx context.Context, message Message) (EmailAnswer, error) {
	answer := EmailAnswer{}

	url := sendUrl
	if message.Mail.IdTemplateMailUser != 0 {
		url = sendByTemplateUrl
	}

	if message.Mail.To.Email == "" || message.Mail.From.Email == "" {
		return answer, &CustomError{Message: "email 'to' and email 'from' are required"}
	}

	if message.Mail.To.Email == message.Mail.From.Email {
		return answer, &CustomError{Message: "email 'to' and email 'from' cannot be the same"}
	}

	if message.Mail.IdTemplateMailUser == 0 && message.Mail.Html == "" && message.Mail.Text == "" {
		return answer, &CustomError{Message: "if no template is used either 'html' or 'text' must be provided"}
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return answer, &CustomError{Message: fmt.Sprintf("marshaling request failed: %s", err.Error())}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+url, bytes.NewReader(payload))
	if err != nil {
		return answer, &CustomError{Message: fmt.Sprintf("creating request failed: %s", err.Error())}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return answer, &CustomError{Message: fmt.Sprintf("sending request failed: %s", err.Error())}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return answer, &CustomError{Message: fmt.Sprintf("reading response failed: %s", err.Error())}
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		if jsonErr := json.Unmarshal(body, &errResp); jsonErr != nil {
			return answer, &CustomError{
				StatusCode: resp.StatusCode,
				Message:    string(body),
			}
		}
		return answer, &CustomError{
			StatusCode: errResp.StatusCode,
			Message:    errResp.Message,
		}
	}

	if err := json.Unmarshal(body, &answer); err != nil {
		return answer, &CustomError{Message: fmt.Sprintf("parsing successful response failed: %s", err.Error())}
	}

	return answer, nil
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("API request failed: [%d] %s", e.StatusCode, e.Message)
}
