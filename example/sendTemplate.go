package main

import (
	"context"
	"fmt"
	"github.com/xjl0/rusender-go"
	"net/http"
)

func main() {
	client := rusender.NewClient(&http.Client{}, "api_key")
	answer, err := client.Send(context.Background(), rusender.Message{
		IdempotencyKey: "12456",
		Mail: rusender.Mail{
			To: rusender.Contact{
				Email: "example@example.com",
				Name:  "example",
			},
			From: rusender.Contact{
				Email: "example2@example.com",
				Name:  "example2",
			},
			Subject: "Test",
			Html:    "<h1>Test</h1>",
		},
	})
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}

	fmt.Printf("answer uuid: %s", answer.Uuid)
}
