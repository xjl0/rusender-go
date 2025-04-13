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
	"github.com/google/uuid"
	"net/http"
	"github.com/xjl0/rusender-go"
)

func main() {
	client := rusender.NewClient(&http.Client{}, "your_api_key")
	answer, err := client.Send(context.Background(), rusender.Message{
		IdempotencyKey: uuid.New().String(),
		Mail: rusender.Mail{
			To: rusender.Contact{
				Email: "example@example.com",
				Name:  "Example",
			},
			From: rusender.Contact{
				Email: "example2@example.com",
				Name:  "Example2",
			},
			Subject:            "Subject",
			IdTemplateMailUser: 1234,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("answer uuid: %s", answer.Uuid)
}
```

## Использование

As is. В первую очередь для личных нужд.