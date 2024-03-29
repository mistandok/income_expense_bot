package telegram

import (
	"encoding/json"
	"errors"
	"income_expense_bot/internal/lib/e"
	"net/url"
	"strconv"
)

func (m *OutgoingMessage) AsUrlValues() (*url.Values, error) {
	q := make(url.Values)
	q.Add("chat_id", strconv.Itoa(m.ChatId))
	q.Add("text", m.Text)

	switch res := m.ReplyMarkup.(type) {
	case InlineKeyboardMarkup, ReplyKeyboardMarkup:
		marchaled, err := json.Marshal(res)
		if err != nil {
			return nil, e.Wrap("can't marshal ReplyMarkup", err)
		}
		q.Add("reply_markup", string(marchaled))
	case nil:
	default:
		return nil, errors.New("undefined interface for ReplyMarkup")
	}

	return &q, nil
}

func (m *OutgoingCallbackMessage) AsUrlValues() (*url.Values, error) {
	q := make(url.Values)
	q.Add("callback_query_id", m.CallbackQueryId)

	if m.Text != nil {
		q.Add("text", *m.Text)
	}

	if m.ShowAlert != nil {
		q.Add("show_alert", strconv.FormatBool(*m.ShowAlert))
	}

	return &q, nil
}
