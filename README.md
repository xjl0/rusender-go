# Неофициальный клиент для API RuSender на Go 

Неофициальный Go-клиент для интеграции с API сервиса [RuSender](https://rusender.ru/) (сервиса отправки email-рассылок). Этот клиент позволяет легко и быстро отправлять письма через API Rusender, используя Go.

## Документация 

https://rusender.ru/developer/api/email/

## Возможности

- Отправка писем через API RuSender
- Поддержка отправки через шаблоны
- Валидация параметров письма перед отправкой
- Обработка ошибок API с детализированными сообщениями
- Поддержка контекста (context.Context)

## Установка

```bash
go get github.com/xjl0/rusender
```

## Пример использования
```go
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

```

## Использование

As is. В первую очередь для личных нужд.