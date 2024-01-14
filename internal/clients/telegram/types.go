package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message"`
	CallbackQuery *CallbackQuery   `json:"callback_query"`
}

type TelegramResponse struct {
	Ok          bool    `json:"ok"`
	ErrorCode   *int    `json:"error_code,omitempty"`
	Description *string `json:"description,omitempty"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type OutgoingMessage struct {
	ChatId      int    `json:"chat_id"`
	Text        string `json:"text"`
	ReplyMarkup any    `json:"reply_markup"` // Доступны типы InlineKeyboardMarkup, ReplyKeyboardMarkup
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type CallbackQuery struct {
	ID      string           `json:"id"`
	From    From             `json:"user"`
	Message *IncomingMessage `json:"message"`
	Data    *string          `json:"data"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string  `json:"text"`
	CallbackData *string `json:"callback_data,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard       [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard *bool              `json:"resize_keyboard,omitempty"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}
