package rusender

// Message https://rusender.ru/developer/api/email/
type Message struct {
	IdempotencyKey string `json:"idempotencyKey,omitempty"` // Ключ идемпотентности используется для предотвращения повторного выполнения одного и того же запроса. По умолчанию, если в течение одного часа поступает запрос на отправку письма с полностью совпадающими параметрами: отправитель (From), получатель (To), тело письма (Body). Такой запрос будет отклонён для предотвращения задвоения писем.
	Mail           Mail   `json:"mail"`                     // Содержит информацию о письме, которое будет отправлено. Поле обязательно для заполнения
}

// Mail Блок содержит информацию о письме, которое будет отправлено
type Mail struct {
	To                 Contact             `json:"to"`                           // Содержит адрес получателя письма
	From               Contact             `json:"from"`                         // Содержит адрес отправителя письма
	Subject            string              `json:"subject"`                      // Содержит тему или заголовок письма
	PreviewTitle       string              `json:"previewTitle,omitempty"`       // Прехедер письма, до 120 символов. https://rusender.ru/glossary/chto-takoe-preheder/
	IdTemplateMailUser int                 `json:"idTemplateMailUser,omitempty"` // Идентификатор шаблона письма, который будет использоваться для отправки письма
	Params             map[string]string   `json:"params,omitempty"`             // Кастомные переменные для вставки в шаблон https://rusender.ru/knowledge-base/email-service/email-templates/peremennye-i-personalizaciya-rassylok/
	Attachments        []map[string]string `json:"attachments,omitempty"`        // Вложение в письмо (файл), в формате массива файлов структурой {"название файла.расширение":"тело файла закодированное в base64"}
	Headers            map[string]string   `json:"headers,omitempty"`            // Системные заголовки письма (необязательно поле, для опытных пользователей) https://nodemailer.com/message/custom-headers
	Html               string              `json:"html,omitempty"`               // Если передать и текстовую и HTML-версию одновременно, то клиент почты получателя будет решать, какую версию отобразить пользователю в зависимости от его настроек и возможностей. Обычно почтовые клиенты отображают в формате HTML, если они поддерживают эту функцию. Наш сервис автоматически генерирует text похожий на html, если text не передан (или передана пустая строка)
	Text               string              `json:"text,omitempty"`               // Если передать и текстовую и HTML-версию одновременно, то клиент почты получателя будет решать, какую версию отобразить пользователю в зависимости от его настроек и возможностей. Обычно почтовые клиенты отображают в формате HTML, если они поддерживают эту функцию. Наш сервис автоматически генерирует text похожий на html, если text не передан (или передана пустая строка)
	Cc                 string              `json:"cc,omitempty"`                 // Это поле «копия» или «отправить копию». Адресат указанный в CC получит копию сообщения, но все получатели смогут видеть, кому еще были отправлены копии сообщения
	Bcc                string              `json:"bcc,omitempty"`                // Это поле «скрытая копия». Это может быть полезно, если вы хотите отправить копию сообщения кому-то без раскрытия его адреса другим адресатам
}

// Contact Описание структуры получателя письма и отправителя
type Contact struct {
	Email string `json:"email"` // Адрес электронной почты получателя. Должен быть действительным адресом электронной почты и не превышать 255 символов
	Name  string `json:"name"`  // Имя получателя. Должно содержать только буквы, цифры и пробелы. Не должно превышать 255 символов
}

type EmailAnswer struct {
	Uuid string `json:"uuid"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type CustomError struct {
	StatusCode int
	Message    string
}
