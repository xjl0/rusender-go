package rusender

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Send_Success(t *testing.T) {
	testCases := []struct {
		result   string
		expected string
		message  Message
	}{
		{
			expected: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			result:   "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			message: Message{
				IdempotencyKey: "42342",
				Mail: Mail{
					To:      Contact{Email: "example@example.com", Name: "Example"},
					From:    Contact{Email: "example2@example.com", Name: "Example2"},
					Subject: "Example Subject",
					Text:    "Example email text",
				},
			},
		}, {
			expected: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			result:   "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			message: Message{
				IdempotencyKey: "34536",
				Mail: Mail{
					To:                 Contact{Email: "example@example.com", Name: "Example"},
					From:               Contact{Email: "example2@example.com", Name: "Example2"},
					Subject:            "Example Subject",
					IdTemplateMailUser: 1234,
				},
			},
		}, {
			expected: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			result:   "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			message: Message{
				IdempotencyKey: "2434656",
				Mail: Mail{
					To:      Contact{Email: "example@example.com", Name: "Example"},
					From:    Contact{Email: "example2@example.com", Name: "Example2"},
					Subject: "Example Subject",
					Html:    "<h1>Example</h1>",
					Text:    "Example",
				},
			},
		}, {
			expected: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			result:   "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			message: Message{
				IdempotencyKey: "ert34364",
				Mail: Mail{
					To:                 Contact{Email: "example@example.com", Name: "Example"},
					From:               Contact{Email: "example2@example.com", Name: "Example2"},
					Subject:            "Example Subject",
					IdTemplateMailUser: 1234,
					Params:             map[string]string{"key": "value", "key2": "value2"},
					Attachments:        []map[string]string{{"file.txt": "VGhpcyBpcyBhIHRlc3Q="}},
				},
			},
		},
	}

	for _, testCase := range testCases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			err := json.NewEncoder(w).Encode(EmailAnswer{Uuid: testCase.expected})
			assert.NoError(t, err)
		}))

		client := NewClient(ts.Client(), "testAPIKey")
		client.baseURL = ts.URL

		resp, err := client.Send(context.Background(), testCase.message)
		assert.NoError(t, err)
		assert.Equal(t, testCase.result, resp.Uuid)
		ts.Close()
	}
}

func TestClient_Auth_Success(t *testing.T) {
	expected := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	testAPIKey := "test-api-key"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedKey := r.Header.Get("X-Api-Key")
		if receivedKey != testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			err := json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid API key"})
			assert.NoError(t, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		err := json.NewEncoder(w).Encode(EmailAnswer{Uuid: expected})
		assert.NoError(t, err)
	}))
	defer ts.Close()

	client := NewClient(ts.Client(), testAPIKey)
	client.baseURL = ts.URL

	msg := Message{
		Mail: Mail{
			To:      Contact{Email: "example@example.com", Name: "Example"},
			From:    Contact{Email: "example2@example.com", Name: "Example2"},
			Subject: "Example Subject",
			Text:    "Example email text",
		},
	}

	resp, err := client.Send(context.Background(), msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp.Uuid)
}

func TestClient_Auth_Invalid(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != "valid-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}))
	defer ts.Close()

	client := NewClient(ts.Client(), "invalid-key")
	client.baseURL = ts.URL

	msg := Message{
		Mail: Mail{
			To:      Contact{Email: "example@example.com", Name: "Example"},
			From:    Contact{Email: "example2@example.com", Name: "Example2"},
			Subject: "Example Subject",
			Text:    "Example email text",
		},
	}
	_, err := client.Send(context.Background(), msg)

	assert.Error(t, err)
	assert.IsType(t, &CustomError{}, err)
	assert.Equal(t, http.StatusUnauthorized, err.(*CustomError).StatusCode)
}

func TestClient_Message_Invalid(t *testing.T) {
	testCases := []struct {
		message Message
	}{
		{
			message: Message{
				Mail: Mail{
					To: Contact{Email: "example@example.com", Name: "Example"},
				},
			},
		}, {
			message: Message{
				Mail: Mail{
					To:   Contact{Email: "example@example.com", Name: "Example"},
					From: Contact{Email: "example@example.com", Name: "Example"},
				},
			},
		}, {
			message: Message{
				Mail: Mail{
					To:   Contact{Email: "example@example.com", Name: "Example"},
					From: Contact{Email: "example2@example.com", Name: "Example2"},
				},
			},
		}, {
			message: Message{
				Mail: Mail{
					To:      Contact{Email: "example@example.com", Name: "Example"},
					From:    Contact{Email: "", Name: "Example2"},
					Subject: "Example Subject",
				},
			},
		},
	}

	for _, testCase := range testCases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))

		client := NewClient(ts.Client(), "testAPIKey")
		client.baseURL = ts.URL

		_, err := client.Send(context.Background(), testCase.message)
		assert.Error(t, err)
		assert.IsType(t, &CustomError{}, err)
		ts.Close()
	}
}

func TestClient_Message_Parse(t *testing.T) {
	testCases := []struct {
		result  string
		message Message
	}{
		{
			result: `{"mail":{"to":{"email":"example@example.com","name":"Example"},"from":{"email":"example2@example.com","name":"Example2"},"subject":"Example Subject","text":"Example email text"}}`,
			message: Message{
				Mail: Mail{
					To:      Contact{Email: "example@example.com", Name: "Example"},
					From:    Contact{Email: "example2@example.com", Name: "Example2"},
					Subject: "Example Subject",
					Text:    "Example email text",
				},
			},
		}, {
			result: `{"mail":{"to":{"email":"example@example.com","name":"Example"},"from":{"email":"example2@example.com","name":"Example2"},"subject":"Example Subject","idTemplateMailUser":1234}}`,
			message: Message{
				Mail: Mail{
					To:                 Contact{Email: "example@example.com", Name: "Example"},
					From:               Contact{Email: "example2@example.com", Name: "Example2"},
					Subject:            "Example Subject",
					IdTemplateMailUser: 1234,
				},
			},
		}, {
			result: `{"idempotencyKey":"3fa85f64-5717-4562-b3fc-2c963f66afa6","mail":{"to":{"email":"example@example.com","name":"Example"},"from":{"email":"example2@example.com","name":"Example2"},"subject":"Example Subject","idTemplateMailUser":1234}}`,
			message: Message{
				IdempotencyKey: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				Mail: Mail{
					To:                 Contact{Email: "example@example.com", Name: "Example"},
					From:               Contact{Email: "example2@example.com", Name: "Example2"},
					Subject:            "Example Subject",
					IdTemplateMailUser: 1234,
				},
			},
		}, {
			result: `{"mail":{"to":{"email":"example@example.com","name":"Example"},"from":{"email":"example2@example.com","name":"Example2"},"subject":"Subject","idTemplateMailUser":1234}}`,
			message: Message{
				Mail: Mail{
					To:                 Contact{Email: "example@example.com", Name: "Example"},
					From:               Contact{Email: "example2@example.com", Name: "Example2"},
					Subject:            "Subject",
					IdTemplateMailUser: 1234,
				},
			},
		}, {
			result: `{"mail":{"to":{"email":"example@example.com","name":"Example"},"from":{"email":"example2@example.com","name":"Example2"},"subject":"Subject","idTemplateMailUser":1234,"params":{"key":"value","key2":"value2"}}}`,
			message: Message{
				Mail: Mail{
					To:                 Contact{Email: "example@example.com", Name: "Example"},
					From:               Contact{Email: "example2@example.com", Name: "Example2"},
					Subject:            "Subject",
					IdTemplateMailUser: 1234,
					Params:             map[string]string{"key": "value", "key2": "value2"},
				},
			},
		}, {
			result: `{"mail":{"to":{"email":"example@example.com","name":"Example"},"from":{"email":"example2@example.com","name":"Example2"},"subject":"Subject","idTemplateMailUser":1234,"params":{"key":"value","key2":"value2"},"attachments":[{"file.pdf":"VGhpcyBpcyBhIHRlc3Q=","file2.pdf":"VGhpcyBpcyBhIHRlc3Q="}]}}`,
			message: Message{
				Mail: Mail{
					To:                 Contact{Email: "example@example.com", Name: "Example"},
					From:               Contact{Email: "example2@example.com", Name: "Example2"},
					Subject:            "Subject",
					IdTemplateMailUser: 1234,
					Params:             map[string]string{"key": "value", "key2": "value2"},
					Attachments:        []map[string]string{{"file.pdf": "VGhpcyBpcyBhIHRlc3Q=", "file2.pdf": "VGhpcyBpcyBhIHRlc3Q="}},
				},
			},
		},
	}

	for _, testCase := range testCases {
		message := Message{}
		assert.NoError(t, json.Unmarshal([]byte(testCase.result), &message))
		assert.Equal(t, testCase.message, message)

		data, err := json.Marshal(&testCase.message)
		assert.NoError(t, err)
		assert.Equal(t, string(data), testCase.result)
	}
}
